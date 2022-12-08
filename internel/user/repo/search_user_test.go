package user

import (
	"context"
	"testing"

	"github.com/exfly/vulcan/pkg/query"
	"github.com/stretchr/testify/require"
)

func TestSQLStore_SearchUser(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	store := NewStore(testDB)

	username := "_test_alice"
	_, err := store.CreateUser(ctx, &CreateUserParams{Username: username})
	if err != nil {
		t.Log(err)
	}

	result, err := store.SearchUser(ctx, &SearchUserParams{
		Username: &query.StringQuery{
			query.SingleStringQuery{
				Op: "=",
				T:  "1",
			},
			query.SingleStringQuery{
				Op: "=",
				T:  "2",
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, username, result.User[0].Username)
}

func TestSQLStore_GetQuery(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	store := NewStore(testDB)

	username := "_test_alice"
	_, err := store.CreateUser(ctx, &CreateUserParams{Username: username})
	if err != nil {
		t.Log(err)
	}

	result, err := store.SearchUser(ctx, &SearchUserParams{
		Username: &query.StringQuery{
			query.SingleStringQuery{
				Op: "=",
				T:  "1",
			},
			query.SingleStringQuery{
				Op: "=",
				T:  "2",
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, username, result.User[0].Username)
}
