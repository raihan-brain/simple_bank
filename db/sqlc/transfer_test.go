package db

import (
	"context"
	"database/sql"
	"github.com/raihan-brain/simple-bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransferParams{
		account1.ID,
		account2.ID,
		util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)

	transfer, err = testQueries.GetTransferByID(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)
}

func TestGetTransferById(t *testing.T) {
	transfer := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransferByID(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer2.Amount, transfer.Amount)
	require.Equal(t, transfer2.ToAccountID, transfer.ToAccountID)
	require.Equal(t, transfer2.FromAccountID, transfer2.FromAccountID)
	require.WithinDuration(t, transfer2.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestListOfTransfer(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListOfTransferParams{
		5,
		5,
	}

	transfers, err := testQueries.ListOfTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, account := range transfers {
		require.NotEmpty(t, account)
	}
}
