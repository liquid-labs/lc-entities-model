package entities

import (
  "github.com/go-pg/pg/orm"
)

func (e *Entity) PrepInsert(q *orm.Query) *orm.Query {
  if e.GetOwnerID() == 0 && e.GetOwnerPubID() != `` {
    q.
      Value(`owner_id`, `(SELECT id FROM entities AS owner_e WHERE owner_e.pub_id=?)`, e.GetOwnerPubID()).
      Returning(`id, pub_id, publicly_readable, created_at, last_updated, deleted_at, owner_id`)
  }
  return q.ExcludeColumn(`owner_pub_id`)
}
