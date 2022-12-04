package user

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type TimeRange struct {
	From sql.NullTime
	To   sql.NullTime
}

func (r TimeRange) Build(field string, builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if r.From.Valid {
		builder = builder.Where(squirrel.GtOrEq{field: r.From.Time})
	}

	if r.To.Valid {
		builder = builder.Where(squirrel.Lt{field: r.To.Time})
	}

	return builder
}

type SearchUserParams struct {
	Username          []string
	FullName          []string
	Email             []string
	PasswordChangedAt *TimeRange `json:"password_changed_at"`
	CreatedAt         *TimeRange `json:"created_at"`
}

func (r SearchUserParams) Build(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if len(r.Username) > 0 {
		builder = builder.Where(squirrel.Eq{"username": r.Username})
	}
	if len(r.FullName) > 0 {
		builder = builder.Where(squirrel.Eq{"fullname": r.FullName})
	}
	if r.PasswordChangedAt != nil {
		builder = r.PasswordChangedAt.Build("password_changed_at", builder)
	}
	if r.CreatedAt != nil {
		builder = r.CreatedAt.Build("created_at", builder)
	}
	return builder
}

type SearchUserResult struct {
	User []User
}

func (store *SQLStore) SearchUser(ctx context.Context, arg *SearchUserParams) (SearchUserResult, error) {
	var result SearchUserResult

	builder := store.db.Builder.Select("*")

	rawSQL, args, err := arg.Build(builder).From("users").ToSql()
	if err != nil {
		return result, err
	}

	err = store.db.SelectContext(ctx, &result.User, rawSQL, args...)
	if err != nil {
		return result, errors.Wrap(err, "")
	}

	return result, err
}
