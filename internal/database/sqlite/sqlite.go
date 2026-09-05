package sqlite

import (
	"context"
	"database/sql"
	"log"

	"github.com/Mrigaank-9/job-scheduler/internal/database"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func SetupSQL(uri string) database.Database {
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		log.Fatal("Unable to create the DB object : ", err)
	}
	return &Sqlite{Db: db}
}

func (s *Sqlite) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.Db.ExecContext(ctx, query, args...)
}

func (s *Sqlite) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.Db.QueryContext(ctx, query, args...)
}

func (s *Sqlite) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.Db.QueryRowContext(ctx, query, args...)
}

func (s *Sqlite) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.Db.PrepareContext(ctx, query)
}

func (s *Sqlite) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.Db.BeginTx(ctx, opts)
}

func (s *Sqlite) Close() error {
	return s.Db.Close()
}
