package entities_test

import (
  "os"
  "testing"
  "time"

  "github.com/go-pg/pg"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/terror/go/terror"
  /* pkg2test */ . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

func init() {
  terror.EchoErrorLog()
}

func retrieveEntity(id EID) (*Entity, terror.Terror) {
  e := &Entity{ID: id}
  if err := rdb.Connect().Model(e).Where(`entity.id=?id`).Select(); err != nil && err != pg.ErrNoRows {
    return nil, terror.ServerError(`Problem retrieving entity.`, err)
  } else if err == pg.ErrNoRows {
    return nil, nil
  } else {
    return e, nil
  }
}

type EntityIntegrationSuite struct {
  suite.Suite
  DB *pg.DB
}
func (s *EntityIntegrationSuite) SetupSuite() {
  s.DB = rdb.Connect()
}
func (s *EntityIntegrationSuite) TearDownSuite() {
  s.DB.Close()
}
func TestEntityIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(EntityIntegrationSuite))
  }
}

func checkDefaults(t *testing.T, e *Entity) {
  assert.NotEqual(t, EID(``), e.GetID(), `ID should have been set on insert.`)
  assert.NotEqual(t, time.Time{}, e.GetCreatedAt(), `'Created at' should have been set on insert.`)
  assert.NotEqual(t, time.Time{}, e.GetLastUpdated(), `'Last updated' should have been set on insert.`)
}

func (s *EntityIntegrationSuite) TestEntityCreateSelfOwner() {
  e1 := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  // model_test verifies that ID, PubID, CreatedAt, LastUpdated, and DeletedAt
  // are initialized to zero/empty values.
  require.NoError(s.T(), rdb.Connect().Insert(e1), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e1)
  assert.Equal(s.T(), e1.GetID(), e1.GetOwnerID())
}

func (s *EntityIntegrationSuite) TestEntityCreateWithOwner() {
  e1 := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), rdb.Connect().Insert(e1))
  e2 := NewEntity(&TestEntity{}, `name`, `description`, e1.GetID(), false)
  // model_test verifies that ID, PubID, CreatedAt, LastUpdated, and DeletedAt
  // are initialized to zero/empty values.

  require.NoError(s.T(), rdb.Connect().Insert(e2), `Unexpected error creating test entity`)
  assert.Equal(s.T(), e1.GetID(), e2.GetOwnerID())
  checkDefaults(s.T(), e2)
}

func (s *EntityIntegrationSuite) TestEntityRetrieve() {
  e := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), rdb.Connect().Insert(e), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityUpdate() {
  e := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), rdb.Connect().Insert(e), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  e.SetName(`foo`)
  e.SetDescription(`bar`)
  e.SetPubliclyReadable(true)
  require.NoError(s.T(), rdb.Connect().Update(e))
  assert.Equal(s.T(), `foo`, e.GetName())
  assert.Equal(s.T(), `bar`, e.GetDescription())
  assert.Equal(s.T(), true, e.IsPubliclyReadable())
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.NotEqual(s.T(), e.GetLastUpdated(), eCopy.GetLastUpdated())
  eCopy.LastUpdated = e.LastUpdated
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityArchive() {
  e := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), rdb.Connect().Insert(e), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  // go-pg v8: last updated and 'deleted_at' get out of sync if not returned. E.g.:
  // rdb.Connect().Delete(e) will fail the test.
  _, err := rdb.Connect().Model(e).Returning(`"last_updated", "deleted_at"`).Where(`entity.id=?id`).Delete()
  require.NoError(s.T(), err)

  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Nil(s.T(), eCopy)

  archived := &Entity{ID: e.GetID()}
  assert.NoError(s.T(), rdb.Connect().Model(archived).Where(`entity.id=?id`).Deleted().Select())
  assert.Equal(s.T(), e, archived)
}
