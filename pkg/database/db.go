package database

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB

	Builder squirrel.StatementBuilderType
}
