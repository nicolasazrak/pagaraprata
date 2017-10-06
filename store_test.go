package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store, err := NewStore("test.db")
	require.NoError(t, err)

	err = store.Clean()
	require.NoError(t, err)

	defer store.Close()

	debt := &Debt{
		ID:     "efg",
		Secret: "abc",
		Data:   &DebtData{"a": "abc"},
	}

	t.Run("get non existen debt", func(t *testing.T) {
		val, err := store.GetDebt(debt.ID)
		require.NoError(t, err)
		require.Nil(t, val)
	})

	t.Run("can save debt", func(t *testing.T) {
		err := store.SaveDebt(debt)
		require.NoError(t, err)
	})

	t.Run("can get existent debt", func(t *testing.T) {
		readDebt, err := store.GetDebt(debt.ID)
		require.NoError(t, err)
		require.Equal(t, *debt, *readDebt)
	})

	t.Run("creates a new debt", func(t *testing.T) {
		initialData := &DebtData{"b": 4}
		newDebt := store.NewDebt(initialData)
		require.NotEmpty(t, newDebt.ID)
		require.NotEmpty(t, newDebt.Secret)
		require.Equal(t, initialData, newDebt.Data)
	})
}
