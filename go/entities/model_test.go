package entities_test

import (
  "reflect"
  "testing"
  "time"

  . "github.com/Liquid-Labs/lc-entity-model/go/entities"
  "github.com/stretchr/testify/assert"
)

func TestEntitiesClone(t *testing.T) {
  now := time.Now()
  orig := NewEntity(`john`, `cool`, `owner-A`, true)
  orig.ID = 1
  orig.PubID = `abc`
  orig.OwnerID = 2
  orig.CreatedAt = now
  orig.LastUpdated = now.Add(100)
  orig.DeletedAt = now.Add(200)
  clone := orig.Clone()

  assert.Equal(t, orig, clone, "Clone does not match.")

  clone.ID = 3
  clone.PubID = `hij`
  clone.Name = `sally`
  clone.Description = `awesome`
  clone.OwnerID = 4
  clone.OwnerPubID = `owner-B`
  clone.PubliclyReadable = false
  clone.CreatedAt = orig.CreatedAt.Add(20)
  clone.LastUpdated = orig.LastUpdated.Add(20)
  clone.DeletedAt = orig.DeletedAt.Add(20)

  // TODO: abstract this
  oReflection := reflect.ValueOf(orig).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  for i := 0; i < oReflection.NumField(); i++ {
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
  orig := NewEntity(`john`, `cool`, `owner-A`, true)
  orig.ID = 1
  orig.PubID = `abc`
  orig.OwnerID = 2
  orig.CreatedAt = now
  orig.LastUpdated = now.Add(100)
  orig.DeletedAt = now.Add(200)
  clone := orig.CloneNew()

  assert.Equal(t, InternalID(0), clone.ID)
  assert.Equal(t, PublicID(``), clone.PubID)
  clone.ID = 1
  clone.PubID = `abc`
  assert.Equal(t, orig, clone, "Clone does not match.")

  clone.ID = 3
  clone.PubID = `hij`
  clone.Name = `sally`
  clone.Description = `awesome`
  clone.OwnerID = 4
  clone.OwnerPubID = `owner-B`
  clone.PubliclyReadable = false
  clone.CreatedAt = orig.CreatedAt.Add(20)
  clone.LastUpdated = orig.LastUpdated.Add(20)
  clone.DeletedAt = orig.DeletedAt.Add(20)

  // TODO: abstract this
  oReflection := reflect.ValueOf(orig).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  for i := 0; i < oReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      oReflection.Field(i).Interface(),
      cReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match.`,
      oReflection.Type().Field(i),
    )
  }
}
