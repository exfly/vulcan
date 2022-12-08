package query

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestStringQuery_Query(t *testing.T) {
	b := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := b.Select("*")
	query := StringQuery{
		{Op: "=", T: "x"},
	}

	builder, err := query.Query("name", builder)
	require.NoError(t, err)

	query = StringQuery{
		{Op: "=", T: "x"},
		{Op: "like", T: "x"},
	}

	builder, err = query.Query("username", builder)
	require.NoError(t, err)

	find := Find()
	find.SetPagination(0, 10)

	builder, err = find.Query("", builder)
	require.NoError(t, err)

	rawQ, args, err := builder.From("users").ToSql()
	require.NoError(t, err)
	spew.Dump(rawQ, args)

	require.Equal(
		t,
		"SELECT * FROM users WHERE name = $1 AND (username = $2 OR username LIKE $3) LIMIT 10",
		rawQ,
	)
}

func TestStringQuery_QueryX(t *testing.T) {
	b := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := b.Select("*")

	queryParams := &struct {
		Name     *StringQuery
		Username *StringQuery
		Password *StringQuery
	}{
		Name: &StringQuery{
			{Op: "=", T: "x"},
		},
		Username: &StringQuery{
			{Op: "=", T: "x"},
			{Op: "like", T: "x"},
		},
	}

	var err error

	builder, err = Build(queryParams, builder)
	require.NoError(t, err)

	rawQ, args, err := builder.From("users").ToSql()
	require.NoError(t, err)
	spew.Dump(rawQ, args)

	require.Equal(
		t,
		"SELECT * FROM users WHERE name = $1 AND (username = $2 OR username LIKE $3)",
		rawQ,
	)
}
