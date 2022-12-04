package database

import (
	"context"
	"database/sql"

	"github.com/exfly/vulcan/internel/config"
	"github.com/exfly/vulcan/pkg/database"

	// https://github.com/jackc/pgx/wiki/Getting-started-with-pgx-through-database-sql
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

var _ DBTX = (*database.DB)(nil)

func New(cfg *config.Config) (*database.DB, error) {
	db, err := sqlx.Connect(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	ret := &database.DB{
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	return ret, nil
}
