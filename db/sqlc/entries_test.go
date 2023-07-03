package db

import (
	"Solvery/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

var user User

func createRandomEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		UserEmail: user.Email,
		Amount:    util.RandomInt(100, 1000),
		Comment:   util.RandomString(6),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.Equal(t, entry.UserEmail, arg.UserEmail)
	require.Equal(t, entry.Amount, arg.Amount)
	require.Equal(t, entry.Comment, arg.Comment)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	user = createRandomUser(t)
	createRandomEntry(t)
}

func TestQueries_ListUserEntries(t *testing.T) {
	user = createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListUserEntriesParams{
		UserEmail: user.Email,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListUserEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestQueries_ListAllEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		user = createRandomUser(t)
		createRandomEntry(t)
	}

	arg := ListAllEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListAllEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
