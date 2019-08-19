package entities

import (
  "fmt"

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

func (e *Entity) UpdateRawQueries(db orm.DB) []*orm.Query {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are auto changed
    Returning(`"last_updated", "deleted_at"`)
  return []*orm.Query{ q }
}

func DoRawUpdate(qs []*orm.Query, db orm.DB) Terror {
  for _, q := range qs {
    if res, err := q.Update(); err != nil {
      alias := string(q.GetModel().Table().Alias)
      return ServerError(fmt.Sprintf(`There was a problem updating the %s record.`, alias), err)
    } else if res.RowsAffected() > 1 {
      alias := string(q.GetModel().Table().Alias)
      return ServerError(fmt.Sprintf(`Unexpected multi-row update while updating %s.`, alias), nil)
    }
  }
  return nil
}

// Update updates an Entity record in the DB. As Entities are logically abstract, one would typically only call this as part of another items update sequence.
func (e *Entity) UpdateRaw(db orm.DB) Terror {
  return DoRawUpdate(e.UpdateRawQueries(db), db)
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
