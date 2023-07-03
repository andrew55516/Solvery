package db

import (
	"Solvery/util"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Name:  util.RandomName(),
		Class: util.RandomClass(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Class, user.Class)
	require.Equal(t, arg.Email, user.Email)
	require.Zero(t, user.Credit)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Class, user2.Class)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Credit, user2.Credit)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestQueries_AddUserCredit(t *testing.T) {
	user1 := createRandomUser(t)
	arg := AddUserCreditParams{
		Amount: util.RandomInt(100, 1000),
		Email:  user1.Email,
	}

	user2, err := testQueries.AddUserCredit(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.Class, user2.Class)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Credit+arg.Amount, user2.Credit)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestQueries_ListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Offset: 5,
		Limit:  5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, u := range users {
		require.NotEmpty(t, u)
	}
}
