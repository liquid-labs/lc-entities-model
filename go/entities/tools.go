package entities

import (
  "fmt"
  "reflect"

  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/terror/go/terror"
)

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

func RunStateQueries(qs []*orm.Query, op *stateOp) Terror {
  for _, q := range qs {
    if res, err := op.f(q); err != nil {
      resourceName := q.GetModel().Table().Name
      return ServerError(fmt.Sprintf(`Error attempting %s on %s .`, op.desc, resourceName), nil)
    } else if q.GetModel().Kind() == reflect.Struct && res.RowsAffected() > 1 {
      // each singular state query should only effect one row (yes?); this may turn out to be too restrictive in some case, but we include it for now as a sanity / data integrety check. In particular, the fear is a rogue query lacking the proper 'Where', and so we want to detect the unexpected and complain so it can be caught in test.
      // go-pg may already be doing a similar check, but the docs are so poor it's unclear and don't want to spend enough time reading code to be sure.
      return ServerError(`Unexpected change to multiple rows.`, nil)
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

func RunRetrieveOp(q *orm.Query, op *retrieveOp) (int, Terror) {
  if count, err := op.f(q); err != nil {
    resourceName := q.GetModel().Table().Name
    return count, ServerError(fmt.Sprintf(`Error attempting %s on %s .`, op.desc, resourceName), nil)
  } else {
    return count, nil
  }
}
