package entities_test

import (
  "testing"
  "time"

  "github.com/go-pg/pg"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  "github.com/Liquid-Labs/env/go/env"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/terror/go/terror"
  /* pkg2test */ . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

func init() {
  terror.EchoErrorLog()
}

func retrieveEntity(id EID) (*Entity, terror.Terror) {
  e := &Entity{}
  q := rdb.Connect().Model(e).Where(`entity.id=?`, id)
  if err := q.Select(); err != nil && err != pg.ErrNoRows {
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
  if env.Get(`SKIP_INTEGRATION`) == `true` {
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
  e1 := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), CreateEntityRaw(e1, rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e1)
  assert.Equal(s.T(), e1.GetID(), e1.GetOwnerID())
}

func (s *EntityIntegrationSuite) TestEntityCreateWithOwner() {
  e1 := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), CreateEntityRaw(e1, rdb.Connect()))

  e2 := NewEntity(`entities`, `name`, `description`, e1.GetID(), false)
  require.NoError(s.T(), CreateEntityRaw(e2, rdb.Connect()), `Unexpected error creating test entity`)

  assert.Equal(s.T(), e1.GetID(), e2.GetOwnerID())
  checkDefaults(s.T(), e2)
}

func (s *EntityIntegrationSuite) TestEntityRetrieve() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), CreateEntityRaw(e, rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityUpdate() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), CreateEntityRaw(e, rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  e.SetName(`foo`)
  e.SetDescription(`bar`)
  e.SetPubliclyReadable(true)
  require.NoError(s.T(), e.UpdateRaw(rdb.Connect()))
  assert.Equal(s.T(), `foo`, e.GetName())
  assert.Equal(s.T(), `bar`, e.GetDescription())
  assert.Equal(s.T(), true, e.IsPubliclyReadable())
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityArchive() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), CreateEntityRaw(e, rdb.Connect()), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  require.NoError(s.T(), e.ArchiveRaw(rdb.Connect()))

  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Nil(s.T(), eCopy)

  archived := &Entity{}
  q := rdb.Connect().Model(archived).Where(`entity.id=?`, e.GetID()).Deleted()
  assert.NoError(s.T(), q.Select())
  assert.Equal(s.T(), e, archived)
}

func (s *EntityIntegrationSuite) TestCreateEntityOnProduction() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  currEnv := env.GetType()
  env.Set(env.DefaultEnvTypeKey, `production`)
  assert.Error(s.T(), CreateEntityRaw(e, rdb.Connect()))
  env.Set(env.DefaultEnvTypeKey, currEnv)
}
