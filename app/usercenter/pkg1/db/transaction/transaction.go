package transaction

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"sync"
)

type Context struct {
	//启用事务标识
	beginTransFlag bool
	rawDb          *pg.DB
	session        orm.DB
	lock           sync.Mutex
}

func (transactionContext *Context) Begin() error {
	transactionContext.lock.Lock()
	defer transactionContext.lock.Unlock()
	transactionContext.beginTransFlag = true
	tx, err := transactionContext.rawDb.Begin()
	if err != nil {
		return err
	}
	transactionContext.session = tx
	return nil
}

func (transactionContext *Context) Commit() error {
	transactionContext.lock.Lock()
	defer transactionContext.lock.Unlock()
	if !transactionContext.beginTransFlag {
		return nil
	}
	if v, ok := transactionContext.session.(*pg.Tx); ok {
		err := v.Commit()
		return err
	}
	return nil
}

func (transactionContext *Context) Rollback() error {
	transactionContext.lock.Lock()
	defer transactionContext.lock.Unlock()
	if !transactionContext.beginTransFlag {
		return nil
	}
	if v, ok := transactionContext.session.(*pg.Tx); ok {
		err := v.Rollback()
		return err
	}
	return nil
}

func (transactionContext *Context) DB() orm.DB {
	return transactionContext.session
}

func NewTransactionContext(db *pg.DB) *Context {
	return &Context{
		rawDb:   db,
		session: db,
	}
}

type Conn interface {
	Begin() error
	Commit() error
	Rollback() error
	DB() orm.DB
}

// UseTrans when beginTrans is true , it will begin a new transaction
// to execute the function, recover when  panic happen
func UseTrans(ctx context.Context,
	db *pg.DB,
	fn func(context.Context, Conn) error, beginTrans bool) (err error) {
	var tx Conn
	tx = NewTransactionContext(db)
	if beginTrans {
		if err = tx.Begin(); err != nil {
			return
		}
	}
	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %w", p, e)
			} else {
				err = fmt.Errorf("recoveer from %#v", p)
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %w", err, e)
			}
		} else {
			err = tx.Commit()
		}
	}()

	return fn(ctx, tx)
}
