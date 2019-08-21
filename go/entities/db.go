package entities

import (
  "github.com/go-pg/pg/orm"
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

func (e *Entity) RetrieveByIDQueries(id EID, db orm.DB) *orm.Query {
  return db.Model(e).Where(`"entity".id=?`, id)
}

func (e *Entity) UpdateQueries(db orm.DB) []*orm.Query {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are auto changed
    Returning(`"last_updated", "deleted_at"`)
  return []*orm.Query{ q }
}

func (e *Entity) ArchiveQueries(db orm.DB) []*orm.Query {
  q := db.Model(e).
    Where(`entity.id=?id`).
    // go-pg doesn't know these are updated by trigger
    Returning(`"last_updated", "deleted_at"`)
  return []*orm.Query{ q }
}

func (e *Entity) DeleteQueries(db orm.DB) []*orm.Query {
  // so another undocumented aspect, go-pg appearently expects you to archive, then delete, because it won't delete something unless archived. Sort of like, "put in trash, take out trash."
  q := db.Model(e).Where(`entity.id=?id`)
  q.GetModel().Table().SoftDeleteField = nil
  return []*orm.Query{ q }
}
