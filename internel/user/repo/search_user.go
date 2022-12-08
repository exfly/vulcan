package user

import (
	"context"

	"github.com/exfly/vulcan/pkg/query"

	"github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

// SearchUserParams
//
//go:generate gomodifytags -file search_user.go -all -add-tags db -w -override -transform snakecase -quiet
type SearchUserParams struct {
	Username          *query.StringQuery `db:"username"`
	FullName          *query.StringQuery `db:"full_name"`
	Email             *query.StringQuery `db:"email"`
	PasswordChangedAt *query.TimeRange   `json:"password_changed_at" db:"password_changed_at"`
	CreatedAt         *query.TimeRange   `json:"created_at" db:"created_at"`
}

func (r SearchUserParams) BuildSQL(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if r.Username != nil {
		builder, _ = r.Username.Query("username", builder)
	}
	if r.FullName != nil {
		builder, _ = r.FullName.Query("fullname", builder)
	}
	if r.PasswordChangedAt != nil {
		builder, _ = r.PasswordChangedAt.Query("password_changed_at", builder)
	}
	if r.CreatedAt != nil {
		builder, _ = r.CreatedAt.Query("created_at", builder)
	}
	return builder
}

type SearchUserResult struct {
	User []User `db:"user"`
}

func (store *SQLStore) SearchUser(ctx context.Context, arg *SearchUserParams) (SearchUserResult, error) {
	var result SearchUserResult

	builder := store.db.Builder.Select("*")

	builder, err := query.Build(arg, builder)
	if err != nil {
		return result, errors.Wrap(err, "")
	}

	rawSQL, args, err := builder.From("users").ToSql()
	if err != nil {
		return result, err
	}

	spew.Dump(rawSQL)
	spew.Dump(args)

	err = store.db.SelectContext(ctx, &result.User, rawSQL, args...)
	if err != nil {
		return result, errors.Wrap(err, "")
	}

	return result, err
}
