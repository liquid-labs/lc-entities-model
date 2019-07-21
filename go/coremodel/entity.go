package coremodel

import (
  "time"

  // "github.com/go-pg/pg"

  // "github.com/Liquid-Labs/go-nullable-mysql/nulls"
  // "github.com/Liquid-Labs/catalyst-sql/go/catsql"
)

// Entity is the basic Catalyst base-type. It is used for "independent" objects
// and data. Any item which is directly retrievable and/or an authorization
// target or subject must embed the Entity type. An Entity should be considered
// an "abstract" type and never created alone, but only as part of creating a
// concrete, final type.
type Entity struct {
  // Note, the ID is for internal use only and may or may not be set depending
  // in the source of the item (client or backend).
  ID               int64       `json:"-" sql:",pk"`
  PubID            string      `json:"pubId" pg:",unique,notnull"`
  Name             string      `json:"name"`
  Description      string      `json:"description"`
  OwnerID          int64       `json:"-" pg:",notnull"`
  OwnerPubID       string      `json:"ownerPubId"`
  PubliclyReadable bool        `json:"publiclyReadable" pg:",notnull"`
  Containers       []Container `json:"containers" pg:"many2many:containers"`
  CreatedAt        time.Time   `json:createdAt`
  LastUpdated      time.Time   `json:"lastUpdated"`
  DeletedAt        time.Time   `pg:",soft_delete"`
}

// catsql.ForQuery(`Entity`, func(q *orm.Query) *orm.Query { return q.Join(`JOIN entities ownerEntitiy ON entity.owner_id=ownerEntity.id`).ColumnExpr(`ownerEntity.pub_id AS owner_pub_id`)})

/*func NewEntity(id int64,
    pubID string,
    name string,
    ownerID int64,
    publiclyReadable bool,
    containers []int64,
    createdAt time.Time,
    lastUpdated time.Time,
    deletedAt time.Time) *Entity {
  return &Entity{id, pubID, name, ownerID, '', publiclyReadable, containers, createdAt, lastUpdated, deletedAt}
}*/

func (e *Entity) Clone() *Entity {
  return &Entity{e.ID, e.PubID, e.Name, e.Description, e.OwnerID, e.OwnerPubID, e.PubliclyReadable, e.Containers, e.CreatedAt, e.LastUpdated, e.DeletedAt}
}

func (e *Entity) GetID() int64 { return e.ID }

func (e *Entity) GetPubID() string { return e.PubID }

func (e *Entity) GetOwnerID() int64 { return e.OwnerID }

func (e *Entity) GetOwnerPubID() string { return e.OwnerPubID }
func (e *Entity) SetOwnerPubID(pid string) { e.OwnerPubID = pid }

func (e *Entity) IsPubliclyReadable() bool { return e.PubliclyReadable }
func (e *Entity) SetPubliclyReadable(r bool) { e.PubliclyReadable = r }

func (e *Entity) GetCreatedAt() time.Time { return e.CreatedAt }

func (e *Entity) GetLastUpdated() time.Time { return e.LastUpdated }

// EntityIface is provided for extension by abstract entity types. In most
// situations you will use the Entity struct directly.
// TODO: confirm no longer necessary.
/*type EntityIface interface {
  GetID() int64

  GetPubID() string

  GetOwnerID() string
  // The owner ID is always set based on the OwnerPubID when storing.

  GetOwnerPubID() string
  SetOwnerPubID(string)

  IsPubliclyReadable() bool
  SetPubliclyReadable(bool)

  GetCreatedAt() time.Time

  GetLastUpdated() time.Time
}*/
