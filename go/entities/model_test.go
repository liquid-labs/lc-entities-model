package entities_test

import (
  "reflect"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"

  // the package being tested
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type TestEntity struct {
  Entity
}

func (te *TestEntity) GetResourceName() ResourceName {
  return ResourceName(`testentities`)
}

func TestNoIDOnCreate(t *testing.T) {
  e := NewEntity(&TestEntity{}, `john`, `cool`, `owner-A`, true)
  assert.Equal(t, EID(``), e.GetID())
}

func TestNoCreatedAtOnCreate(t *testing.T) {
  e := NewEntity(&TestEntity{}, `john`, `cool`, `owner-A`, true)
  assert.Equal(t, time.Time{}, e.GetCreatedAt())
}

func TestNoLastUpdatedOnCreate(t *testing.T) {
  e := NewEntity(&TestEntity{}, `john`, `cool`, `owner-A`, true)
  assert.Equal(t, time.Time{}, e.GetLastUpdated())
}

func TestNoDeletedAtOnCreate(t *testing.T) {
  e := NewEntity(&TestEntity{}, `john`, `cool`, `owner-A`, true)
  assert.Equal(t, time.Time{}, e.GetDeletedAt())
}

func TestEntitiesClone(t *testing.T) {
  now := time.Now()
  orig := NewEntity(&TestEntity{}, `john`, `cool`, `owner-A`, true)
  orig.ID = EID(`abc`)
  orig.OwnerID = EID(`owner-A`)
  orig.CreatedAt = now
  orig.LastUpdated = now.Add(100)
  orig.DeletedAt = now.Add(200)
  clone := orig.Clone()

  assert.Equal(t, orig, clone, "Clone does not match.")

  clone.ID = `hij`
  clone.ResourceName = `foos`
  clone.Name = `sally`
  clone.Description = `awesome`
  clone.OwnerID = EID(`owner-B`)
  clone.PubliclyReadable = false
  clone.CreatedAt = orig.CreatedAt.Add(20)
  clone.LastUpdated = orig.LastUpdated.Add(20)
  clone.DeletedAt = orig.DeletedAt.Add(20)

  // TODO: abstract this
  oReflection := reflect.ValueOf(orig).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  // start at '1' to skip the 'tableName' field
  for i := 1; i < oReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      oReflection.Field(i).Interface(),
      cReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match.`,
      oReflection.Type().Field(i),
    )
  }
}

func TestEntitiesCloneNew(t *testing.T) {
  now := time.Now()
  orig := NewEntity(&TestEntity{}, `john`, `cool`, `owner-A`, true)
  orig.ID = `abc`
  orig.OwnerID = `owner-A`
  orig.CreatedAt = now
  orig.LastUpdated = now.Add(100)
  orig.DeletedAt = now.Add(200)
  clone := orig.CloneNew()

  assert.Equal(t, EID(``), clone.ID)
  assert.Equal(t, time.Time{}, clone.CreatedAt)
  assert.Equal(t, time.Time{}, clone.LastUpdated)
  assert.Equal(t, time.Time{}, clone.DeletedAt)

  clone.Name = `sally`
  clone.ResourceName = `foos`
  clone.Description = `awesome`
  clone.OwnerID = `owner-B`
  clone.PubliclyReadable = false

  // TODO: abstract this
  oReflection := reflect.ValueOf(orig).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  // start at '1' to skip the 'tableName' field
  for i := 1; i < oReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      oReflection.Field(i).Interface(),
      cReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match.`,
      oReflection.Type().Field(i),
    )
  }
}
