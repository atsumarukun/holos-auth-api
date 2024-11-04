package infrastructure

import (
	"context"
	"database/sql"
	"holos-auth-api/internal/app/api/domain"
	"log"

	"github.com/jmoiron/sqlx"
)

type transactionKey struct{}

type sqlxTransactionObject struct {
	db *sqlx.DB
}

func NewSqlxTransactionObject(db *sqlx.DB) domain.TransactionObject {
	return &sqlxTransactionObject{
		db: db,
	}
}

func (sto *sqlxTransactionObject) Transaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := sto.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	ctx = context.WithValue(ctx, transactionKey{}, tx)

	if err := fn(ctx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(err.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err.Error())
	}

	return nil
}

type sqlxDriver interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

func getSqlxDriver(ctx context.Context, db *sqlx.DB) sqlxDriver {
	if tx, ok := ctx.Value(transactionKey{}).(*sqlx.Tx); !ok {
		return db
	} else {
		return tx
	}
}
