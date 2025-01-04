package database

import (
	"context"
	"database/sql"
	"holos-auth-api/internal/app/api/domain"
	"log"

	"github.com/jmoiron/sqlx"
)

type transactionKey struct{}

type transactionObject struct {
	db *sqlx.DB
}

func NewDBTransactionObject(db *sqlx.DB) domain.TransactionObject {
	return &transactionObject{
		db: db,
	}
}

func (to *transactionObject) Transaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := to.db.Beginx()
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

type driver interface {
	Rebind(string) string
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
}

func getDriver(ctx context.Context, db *sqlx.DB) driver {
	if tx, ok := ctx.Value(transactionKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
