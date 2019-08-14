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
  /* pkg2test */ "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

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

func checkDefaults(t *testing.T, e *entities.Entity) {
  assert.NotEqual(t, entities.EID(``), e.GetID(), `ID should have been set on insert.`)
  assert.NotEqual(t, time.Time{}, e.GetCreatedAt(), `'Created at' should have been set on insert.`)
  assert.NotEqual(t, time.Time{}, e.GetLastUpdated(), `'Last updated' should have been set on insert.`)
}

func (s *EntityIntegrationSuite) TestEntityInsertNoOwner() {
  e1 := entities.NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  // model_test verifies that ID, PubID, CreatedAt, LastUpdated, and DeletedAt
  // are initialized to zero/empty values.
  _, err := s.DB.Model(e1).Insert()
  require.NoError(s.T(), err, `Unexpected error creating test entity`)
  checkDefaults(s.T(), e1)
}

func (s *EntityIntegrationSuite) TestEntityInsertWithOwner() {
  e1 := entities.NewEntity(&TestEntity{}, `name`, `description`, ``, false)
  _, err := s.DB.Model(e1).Insert()
  e2 := entities.NewEntity(&TestEntity{}, `name`, `description`, e1.GetID(), false)
  // model_test verifies that ID, PubID, CreatedAt, LastUpdated, and DeletedAt
  // are initialized to zero/empty values.

  _, err = s.DB.Model(e2).Insert()
  require.NoError(s.T(), err, `Unexpected error creating test entity`)
  assert.Equal(s.T(), e1.GetID(), e2.GetOwnerID())
  checkDefaults(s.T(), e2)
}
