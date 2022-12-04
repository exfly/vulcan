package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLStore_SearchUser(t *testing.T) {
	ctx := context.Background()
	store := NewStore(testDB)

	username := "_test_alice"
	_, err := store.CreateUser(ctx, &CreateUserParams{Username: username})
	if err != nil {
		t.Log(err)
	}

	result, err := store.SearchUser(ctx, &SearchUserParams{
		Username: []string{username},
	})
	require.NoError(t, err)

	require.Equal(t, username, result.User[0].Username)
}
