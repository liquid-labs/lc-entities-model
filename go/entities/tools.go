package entities

import (
  "fmt"
  "reflect"

  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
)

// retrieval is already easy in the simple case and arbitrary in other cases, so we don't attempt to abstract it. Rather, it's best to just interact with the go-pg tools directly. The item manager deals with transactional, state change methods.
type creatable interface { CreateQueries(orm.DB) []*orm.Query }
type updatable interface { UpdateQueries(orm.DB) []*orm.Query }
type archivable interface { ArchiveQueries(orm.DB) []*orm.Query }
type deletable interface { DeleteQueries(orm.DB) []*orm.Query }

type ItemManager struct {
  db                     *pg.DB
  tx                     *pg.Tx
  AllowUnsafeStateChange bool
}
func NewItemManager(db *pg.DB) *ItemManager {
  return &ItemManager{db:db, AllowUnsafeStateChange:false}
}
func (im *ItemManager) getDB() orm.DB {
  if im.tx != nil { return im.tx } else { return im.db }
}
func (im *ItemManager) StartTransaction() error {
  if im.tx != nil {
    return fmt.Errorf(`Attempt to start transaction while in a transaction.`)
  }
  tx, err := im.db.Begin()
  if err != nil { im.tx = tx }
  return err
}
func (im *ItemManager) dropTransaction() { im.tx = nil }
func (im *ItemManager) CommitTransaction() error {
  if im.tx == nil {
    return fmt.Errorf(`Attempt to commit non-existent transaction.`)
  }
  defer im.dropTransaction()
  return im.tx.Commit()
}
func (im *ItemManager) RollbackTransaction() error {
  if im.tx == nil {
    return fmt.Errorf(`Attempt to rollback non-existent transaction.`)
  }
  defer im.dropTransaction()
  return im.tx.Rollback()
}

func (im *ItemManager) doStateChangeOp(qs []*orm.Query, op *stateOp) error {
  if !im.AllowUnsafeStateChange {
    return fmt.Errorf(`Attempt to perform '%s' outside of transaction context.`, op.desc)
  } else {
    return RunStateQueries(qs, op)
  }
}

func (im *ItemManager) CreateRaw(item creatable) error {
  return im.doStateChangeOp(item.CreateQueries(im.db), CreateOp)
}

func (im *ItemManager) UpdateRaw(item updatable) error {
  return im.doStateChangeOp(item.UpdateQueries(im.db), UpdateOp)
}

func (im *ItemManager) ArchiveRaw(item archivable) error {
  return im.doStateChangeOp(item.ArchiveQueries(im.db), ArchiveOp)
}

func (im *ItemManager) DeleteRaw(item deletable) error {
  return im.doStateChangeOp(item.DeleteQueries(im.db), DeleteOp)
}

// State ops:
// * potentially have mulitple steps (== queries).
// * should, at most, effect 1 row per step.
// * should be run in a transaction (outside test)

// Retrieve ops:
// * always have a single query
// * single retrievals may allow 0 or require 1
// * lists may return any number

type stateFunc func (*orm.Query) (orm.Result, error)

type stateOp struct {
  f    stateFunc
  desc string
}

var CreateOp = &stateOp {
  func (q *orm.Query) (orm.Result, error) { return q.Insert() },
  `create`,
}

var UpdateOp = &stateOp{
  func (q *orm.Query) (orm.Result, error) { return q.Update() },
  `update`,
}

var ArchiveOp = &stateOp{
  func (q *orm.Query) (orm.Result, error) { return q.Delete() },
  `archive`,
}

var DeleteOp = &stateOp{
  func (q *orm.Query) (orm.Result, error) { return q.ForceDelete() },
  `delete`,
}

func RunStateQueries(qs []*orm.Query, op *stateOp) error {
  for _, q := range qs {
    if res, err := op.f(q); err != nil {
      resourceName := q.GetModel().Table().Name
      return fmt.Errorf(`Error attempting %s %s item; %s.`, op.desc, resourceName, err)
    } else if q.GetModel().Kind() == reflect.Struct && res.RowsAffected() > 1 {
      // each singular state query should only effect one row (yes?); this may turn out to be too restrictive in some case, but we include it for now as a sanity / data integrety check. In particular, the fear is a rogue query lacking the proper 'Where', and so we want to detect the unexpected and complain so it can be caught in test.
      // go-pg may already be doing a similar check, but the docs are so poor it's unclear and don't want to spend enough time reading code to be sure.
      return fmt.Errorf(`Unexpected change to multiple rows.`)
    }
  }
  // if we fall out of the loop with no errors, then we're all good.
  return nil
}

type retrieveFunc func (*orm.Query) (int, error)

type retrieveOp struct {
  f    retrieveFunc
  desc string
}

var RetrieveOp = &retrieveOp{
  func (q *orm.Query) (int, error) {
    if err := q.Select(); err != nil {
      if err == pg.ErrNoRows {
        return 0, nil
      } else {
        return -1, err
      }
    }
    return 1, nil
  },
  `retrieve`,
}

var MustRetrieveOp = &retrieveOp{
  func (q *orm.Query) (int, error) {
    if err := q.Select(); err != nil {
      if err == pg.ErrNoRows {
        return 0, err
      } else {
        return -1, err
      }
    }
    return 1, nil
  },
  `must retrieve`,
}

var ListOp = &retrieveOp{
  func (q *orm.Query) (int, error) { return  q.SelectAndCount() },
  `list`,
}

func RunRetrieveOp(q *orm.Query, op *retrieveOp) (int, error) {
  if count, err := op.f(q); err != nil {
    resourceName := q.GetModel().Table().Name
    return count, fmt.Errorf(`Error attempting %s on %s .`, op.desc, resourceName)
  } else {
    return count, nil
  }
}
