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
  e, q := ModelEntity(rdb.Connect())
  if err := q.Where(`entity.id=?`, id).Select(); err != nil && err != pg.ErrNoRows {
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
func (s *EntityIntegrationSuite) TearDownSuite() {
  rdb.Connect().Close()
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
  require.NoError(s.T(), e1.Create(rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e1)
  assert.Equal(s.T(), e1.GetID(), e1.GetOwnerID())
}

func (s *EntityIntegrationSuite) TestEntityCreateWithOwner() {
  e1 := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), e1.Create(rdb.Connect()))

  e2 := NewEntity(&TestEntity{}, `name`, `description`, e1.GetID(), false)
  require.NoError(s.T(), e2.Create(rdb.Connect()), `Unexpected error creating test entity`)

  assert.Equal(s.T(), e1.GetID(), e2.GetOwnerID())
  checkDefaults(s.T(), e2)
}

func (s *EntityIntegrationSuite) TestEntityRetrieve() {
  e := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), e.Create(rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityUpdate() {
  e := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), e.Create(rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  e.SetName(`foo`)
  e.SetDescription(`bar`)
  e.SetPubliclyReadable(true)
  require.NoError(s.T(), e.Update(rdb.Connect()))
  assert.Equal(s.T(), `foo`, e.GetName())
  assert.Equal(s.T(), `bar`, e.GetDescription())
  assert.Equal(s.T(), true, e.IsPubliclyReadable())
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityArchive() {
  e := NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  require.NoError(s.T(), e.Create(rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  require.NoError(s.T(), e.Archive(rdb.Connect()))

  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Nil(s.T(), eCopy)

  archived, q := ModelEntity(rdb.Connect())
  assert.NoError(s.T(), q.Where(`entity.id=?`, e.GetID()).Deleted().Select())
  assert.Equal(s.T(), e, archived)
}
