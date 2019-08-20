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

func (e *Entity) CreateQueries(db orm.DB) []*orm.Query {
  return []*orm.Query{ db.Model(e) }
}

// Create creates (or inserts) a new Entity record into the DB. As Entities are logically abstract, one would typically only call this as part of another items create sequence.
func CreateEntityRaw(eb Entable, db orm.DB) Terror {
  if !eb.IsConcrete() && env.IsProduction() {
    // TODO: improve this error message
    return BadRequestError(`Attempt to create non-concrete entity in prdouction.`)
  }
  return RunStateQueries(eb.GetEntity().CreateQueries(db), CreateOp)
}

func (e *Entity) UpdateQueries(db orm.DB) []*orm.Query {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are auto changed
    Returning(`"last_updated", "deleted_at"`)
  return []*orm.Query{ q }
}

// Update updates an Entity record in the DB. As Entities are logically abstract, one would typically only call this as part of another items update sequence.
func (e *Entity) UpdateRaw(db orm.DB) Terror {
  return RunStateQueries(e.UpdateQueries(db), UpdateOp)
}

func (e *Entity) ArchiveQueries(db orm.DB) []*orm.Query {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are updated by trigger
    Returning(`"last_updated", "deleted_at"`)
  return []*orm.Query{ q }
}

// Archive updates an Entity record in the DB. As Entities are logically abstract, one would typically only call this as part of another items archive sequence.
func (e *Entity) ArchiveRaw(db orm.DB) Terror {
  return RunStateQueries(e.ArchiveQueries(db), ArchiveOp)
}

func (e *Entity) DeleteQueries(db orm.DB) []*orm.Query {
  q := db.Model(e).
    Where(`entity.id=?id`)
  return []*orm.Query{ q }
}

func (e *Entity) Delete(db orm.DB) Terror {
  return RunStateQueries(e.DeleteQueries(db), DeleteOp)
}
