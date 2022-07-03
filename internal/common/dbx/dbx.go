package dbx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"

	"github.com/jmoiron/sqlx"
)

var (
	ErrTransactionStarted = errors.New("transaction has already been started")
	ErrDBType             = errors.New("wrong type of DB interface")
	ErrDBNotFoundRecord   = gorm.ErrRecordNotFound.Error()
)

type Jsonb driver.Value

// DB interface is implemented by sqlx.DB and sqlx.Tx
// so it can be used either by db or transaction instance
type DB interface {
	// GetContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// SelectContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
	Rebind(query string) string
	// BindNamed binds a query using the DB driver's bindvar type.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	// NamedExecContext using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// QueryxContext queries the database and returns an *sqlx.Rows.
	// Any placeholder parameters are replaced with supplied args.
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

func JsonbToBytes(jsonb Jsonb) []byte {
	return jsonb.([]uint8)
}

type txKey struct{}

// withTx assign Tx to the context
func withTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// Connection get connection (*sqlx.db | *sqlx.tx) from the context
func Connection(ctx context.Context, db DB) DB {
	tx, ok := txFromContext(ctx)
	if ok {
		return tx
	}
	return db
}

func txFromContext(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sqlx.Tx)
	return tx, ok
}

// Transactional run wrapped func in transaction
func Transactional(ctx context.Context, db *sqlx.DB, wrappedFunc func(ctx context.Context) error) (err error) {
	tx, ctxWithTx, err := BeginTx(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = wrappedFunc(ctxWithTx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// BeginTx start transaction
// Error ErrTransactionStarted
// Error ErrDBType
func BeginTx(ctx context.Context, db DB) (*sqlx.Tx, context.Context, error) {
	tx, ok := txFromContext(ctx)
	if ok {
		return tx, ctx, ErrTransactionStarted
	}

	dbt, ok := db.(*sqlx.DB)
	if !ok {
		return nil, ctx, ErrDBType
	}

	tx, err := dbt.BeginTxx(ctx, nil)
	if err != nil {
		return nil, ctx, err
	}

	return tx, withTx(ctx, tx), nil
}
