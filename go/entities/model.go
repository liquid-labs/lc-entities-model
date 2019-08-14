package entities

import (
  "time"

  // "github.com/go-pg/pg"
)

type InternalID int64
type PublicID string
type ResourceName string

// Entity is the base type for all independent entities in the Liquid Code
// model. Any item which is directly retrievable, an authorization target, or
// authorization subject must embed the Entity type. An Entity should be
// considered an "abstract" type and never created directly, but only as part of
// a concrete, final type.
type Entity struct {
  tableName        struct{}     `sql:"entities,alias:entity,select:entities_owner_pub_id"`
  // Note, the ID is for internal use only and may or may not be set depending
  // in the source of the item (client or backend).
  ID               InternalID   `json:"-" sql:",pk"`
  PubID            PublicID     `json:"pubId" pg:",unique,notnull"`
  ResourceName     ResourceName `json:"resourceName"`
  Name             string       `json:"name"`
  Description      string       `json:"description"`
  OwnerID          InternalID   `json:"-"`
  OwnerPubID       PublicID     `json:"ownerPubId"`
  PubliclyReadable bool         `json:"publiclyReadable" pg:",notnull"`
  CreatedAt        time.Time    `json:createdAt`
  LastUpdated      time.Time    `json:"lastUpdated"`
  DeletedAt        time.Time    `pg:",soft_delete"`
}

// catsql.ForQuery(`Entity`, func(q *orm.Query) *orm.Query { return q.Join(`JOIN entities ownerEntitiy ON entity.owner_id=ownerEntity.id`).ColumnExpr(`ownerEntity.pub_id AS owner_pub_id`)})

func NewEntity(
    exemplar Identifiable,
    name string,
    description string,
    ownerPubID PublicID,
    publiclyReadable bool) *Entity {
  return &Entity{
    struct{}{},
    exemplar.GetID(),
    exemplar.GetPubID(),
    exemplar.GetResourceName(),
    name,
    description,
    0,
    ownerPubID,
    publiclyReadable,
    time.Time{},
    time.Time{},
    time.Time{},
  }
}

func (e *Entity) Clone() *Entity {
  return &Entity{
    struct{}{},
    e.ID,
    e.PubID,
    e.ResourceName,
    e.Name,
    e.Description,
    e.OwnerID,
    e.OwnerPubID,
    e.PubliclyReadable,
    e.CreatedAt,
    e.LastUpdated,
    e.DeletedAt,
  }
}

func (e *Entity) CloneNew() *Entity {
  newE := e.Clone()
  newE.ID = 0
  newE.PubID = ``
  newE.CreatedAt = time.Time{}
  newE.LastUpdated = time.Time{}
  newE.DeletedAt = time.Time{}
  return newE
}

func (e *Entity) GetID() InternalID { return e.ID }

func (e *Entity) GetPubID() PublicID { return e.PubID }

func (e *Entity) GetResourceName() ResourceName { return e.ResourceName }

func (e *Entity) GetOwnerID() InternalID { return e.OwnerID }

func (e *Entity) GetOwnerPubID() PublicID { return e.OwnerPubID }
func (e *Entity) SetOwnerPubID(pid PublicID) { e.OwnerPubID = pid }

func (e *Entity) IsPubliclyReadable() bool { return e.PubliclyReadable }
func (e *Entity) SetPubliclyReadable(r bool) { e.PubliclyReadable = r }

func (e *Entity) GetCreatedAt() time.Time { return e.CreatedAt }

func (e *Entity) GetLastUpdated() time.Time { return e.LastUpdated }

func (e *Entity) GetDeletedAt() time.Time { return e.DeletedAt }

type Identifiable interface {
  GetID() InternalID
  GetPubID() PublicID
  GetResourceName() ResourceName
}
