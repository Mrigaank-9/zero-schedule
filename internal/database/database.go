package database

import (
	"context"
	"database/sql"
)

type Database interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)

	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	Close() error
}
