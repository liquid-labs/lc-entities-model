package entities

import (
  "github.com/go-pg/pg/orm"

  "github.com/Liquid-Labs/env/go/env"
  . "github.com/Liquid-Labs/terror/go/terror"
)

var EntityFields = []string{
  `resource_name`,
  `name`,
  `description`,
  `owner_id`,
  `publicly_readable`,
  `created_at`,
  `last_updated`,
  `deleted_at`,
}

// Create creates (or inserts) a new Entity record into the DB. As Entities are logically abstract, one would typically only call this as part of another items create sequence.
func CreateEntityRaw(eb Entable, db orm.DB) Terror {
  if !eb.IsConcrete() && env.IsProduction() {
    // TODO: improve this error message
    return BadRequestError(`Attempt to create non-concrete entity in prdouction.`)
  } else {
    if err := db.Insert(eb.GetEntity()); err != nil {
      return ServerError(`There was a problem creating the entity record.`, err)
    } else {
      return nil
    }
  }
}

// Update updates an Entity record in the DB. As Entities are logically abstract, one would typically only call this as part of another items update sequence.
func (e *Entity) UpdateRaw(db orm.DB) Terror {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are auto changed
    Returning(`"last_updated", "deleted_at"`)
  if _, err := q.Update(); err != nil {
    return ServerError(`There was a problem updating the entity record.`, err)
  } else {
    return nil
  }
}

// Archive updates an Entity record in the DB. As Entities are logically abstract, one would typically only call this as part of another items archive sequence.
func (e *Entity) ArchiveRaw(db orm.DB) Terror {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are
    Returning(`"last_updated", "deleted_at"`)
  if _, err := q.Delete(); err != nil {
    return ServerError(`There was a problem deleting the entity record.`, err)
  } else {
    return nil
  }
}
