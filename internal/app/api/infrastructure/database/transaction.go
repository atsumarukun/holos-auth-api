package database

import (
	"context"
	"database/sql"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/pkg/status"
	"log"
	"net/http"

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

func (o *sqlxTransactionObject) Transaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := o.db.Beginx()
	if err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
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
	Rebind(string) string
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row
}

func getSqlxDriver(ctx context.Context, db *sqlx.DB) sqlxDriver {
	if tx, ok := ctx.Value(transactionKey{}).(*sqlx.Tx); !ok {
		return db
	} else {
		return tx
	}
}
