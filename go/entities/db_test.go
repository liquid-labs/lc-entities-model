package entities_test

import (
  "testing"
  "time"

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

func retrieveEntity(id EID) (*Entity, error) {
  e := &Entity{}
  q := rdb.Connect().Model(e).Where(`entity.id=?`, id)
  count, err := RunRetrieveOp(q, RetrieveOp)
  if count == 0 { return nil, err } else { return e, err }
}

type EntityIntegrationSuite struct {
  suite.Suite
  IM *ItemManager
}
func (s *EntityIntegrationSuite) SetupSuite() {
  s.IM = NewItemManager(rdb.Connect())
  s.IM.AllowUnsafeStateChange = true
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
  require.NoError(s.T(), s.IM.CreateRaw(e1), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e1)
  assert.Equal(s.T(), e1.GetID(), e1.GetOwnerID())
}

func (s *EntityIntegrationSuite) TestEntityCreateWithOwner() {
  e1 := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), s.IM.CreateRaw(e1))

  e2 := NewEntity(`entities`, `name`, `description`, e1.GetID(), false)
  require.NoError(s.T(), s.IM.CreateRaw(e2), `Unexpected error creating test entity`)

  assert.Equal(s.T(), e1.GetID(), e2.GetOwnerID())
  checkDefaults(s.T(), e2)
}

func (s *EntityIntegrationSuite) TestEntityRetrieve() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), s.IM.CreateRaw(e), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityUpdate() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), s.IM.CreateRaw(e), `Unexpected error creating test entity`)
  checkDefaults(s.T(), e)
  e.SetName(`foo`)
  e.SetDescription(`bar`)
  e.SetPubliclyReadable(true)
  require.NoError(s.T(), s.IM.UpdateRaw(e))
  assert.Equal(s.T(), `foo`, e.GetName())
  assert.Equal(s.T(), `bar`, e.GetDescription())
  assert.Equal(s.T(), true, e.IsPubliclyReadable())
  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), e, eCopy)
}

func (s *EntityIntegrationSuite) TestEntityArchive() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), s.IM.CreateRaw(e))
  checkDefaults(s.T(), e)
  require.NoError(s.T(), s.IM.ArchiveRaw(e))
  assert.NotEmpty(s.T(), e.GetDeletedAt())

  eCopy, err := retrieveEntity(e.GetID())
  require.NoError(s.T(), err)
  assert.Nil(s.T(), eCopy)

  archived := &Entity{}
  q := rdb.Connect().Model(archived).Where(`entity.id=?`, e.GetID()).Deleted()
  assert.NoError(s.T(), q.Select())
  assert.Equal(s.T(), e, archived)
}

func (s *EntityIntegrationSuite) TestEntityDelete() {
  e := NewEntity(`entities`, `name`, `description`, ``, false)
  require.NoError(s.T(), s.IM.CreateRaw(e))

  var e1, e2 int
  rdb.Connect().Query(&e1, "SELECT COUNT(*) FROM entities")
  require.NoError(s.T(), s.IM.DeleteRaw(e))
  rdb.Connect().Query(&e2, "SELECT COUNT(*) FROM entities")
  assert.Equal(s.T(), e1 - 1, e2)
}
