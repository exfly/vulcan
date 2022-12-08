package query

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type TimeRange struct {
	From sql.NullTime `db:"from"`
	To   sql.NullTime `db:"to"`
}

func (r TimeRange) Query(field string, builder squirrel.SelectBuilder) (squirrel.SelectBuilder, error) {
	if r.From.Valid {
		builder = builder.Where(squirrel.GtOrEq{field: r.From.Time})
	}

	if r.To.Valid {
		builder = builder.Where(squirrel.Lt{field: r.To.Time})
	}

	return builder, nil
}
