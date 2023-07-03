package db

import (
	"Solvery/util"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPSQLStore_PaymentTx(t *testing.T) {
	store := NewStore(testDB)

	user := createRandomUser(t)
	fmt.Println(">> before:", user.Credit)

	// run n concurrent transfer transactions
	n := 5
	amount := int32(10)

	errs := make(chan error)
	results := make(chan PaymentTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.PaymentTx(context.Background(), PaymentTxParams{
				UserEmail: user.Email,
				Amount:    amount,
				Comment:   util.RandomString(6),
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		updUser := result.User
		require.NotEmpty(t, updUser)
		require.Equal(t, user.Name, updUser.Name)
		require.Equal(t, user.Class, updUser.Class)
		require.Equal(t, user.Email, updUser.Email)
		require.WithinDuration(t, updUser.CreatedAt, user.CreatedAt, time.Second)

		diff := updUser.Credit - user.Credit
		require.True(t, diff%amount == 0)

		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		entry := result.Entry
		require.NotEmpty(t, entry)
		require.Equal(t, user.Email, entry.UserEmail)
		require.Equal(t, amount, entry.Amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)

		fmt.Println(">> tx:", updUser.Credit)
	}

	updUser, err := testQueries.GetUser(context.Background(), user.Email)
	require.NoError(t, err)

	fmt.Println(">> after:", updUser.Credit)
	require.Equal(t, user.Credit+int32(n)*amount, updUser.Credit)
}
