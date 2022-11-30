package pgstore

import (
	"context"
	"database/sql"
)

type DestScanner interface {
	Scan(dest ...interface{}) error
}

type DBI interface {
	Close() error
	Exec(ctx context.Context, stmt string, args ...interface{}) (interface{}, error)
	QueryRow(ctx context.Context, stmt string, args ...interface{}) DestScanner
}

type DatabaseSQLAdaptor struct {
	db *sql.DB
}

func NewDatabaseSQLAdaptor(db *sql.DB) *DatabaseSQLAdaptor {
	return &DatabaseSQLAdaptor{db: db}
}

func (a *DatabaseSQLAdaptor) Close() error {
	return a.db.Close()
}

func (a *DatabaseSQLAdaptor) Exec(ctx context.Context, stmt string, args ...interface{}) (interface{}, error) {
	return a.db.ExecContext(ctx, stmt, args...)
}

func (a *DatabaseSQLAdaptor) QueryRow(ctx context.Context, stmt string, args ...interface{}) DestScanner {
	return a.db.QueryRowContext(ctx, stmt, args...)
}
