package infrastructure

import (
	"context"
	"database/sql"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/pkg/apierr"
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

func (o *sqlxTransactionObject) Transaction(ctx context.Context, fn func(context.Context) apierr.ApiError) apierr.ApiError {
	tx, err := o.db.Beginx()
	if err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
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
