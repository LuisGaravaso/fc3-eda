package entity_test

import (
	"testing"

	"balance/internal/entity"

	"github.com/stretchr/testify/assert"
)

func TestNewBalance(t *testing.T) {
	t.Run("should create a new balance", func(t *testing.T) {
		balance, err := entity.NewBalance("account1", 100.0)

		assert.Nil(t, err)
		assert.NotNil(t, balance)
		assert.Equal(t, "account1", balance.AccountId)
		assert.Equal(t, 100.0, balance.Balance)
	})

	t.Run("should return error when account id is empty", func(t *testing.T) {
		balance, err := entity.NewBalance("", 100.0)

		assert.NotNil(t, err)
		assert.Nil(t, balance)
		assert.Equal(t, entity.ErrInvalidClient, err.Error())
	})

	t.Run("should return error when balance is negative", func(t *testing.T) {
		balance, err := entity.NewBalance("account1", -100.0)

		assert.NotNil(t, err)
		assert.Nil(t, balance)
		assert.Equal(t, entity.ErrInsufficientBalance, err.Error())
	})

	t.Run("should return update balance", func(t *testing.T) {
		balance, err := entity.NewBalance("account1", 100.0)
		assert.Nil(t, err)

		err = balance.UpdateBalance(200.0)
		assert.Nil(t, err)
		assert.Equal(t, 200.0, balance.Balance)
	})

	t.Run("should return error when updating balance to negative", func(t *testing.T) {
		balance, err := entity.NewBalance("account1", 100.0)
		assert.Nil(t, err)

		err = balance.UpdateBalance(-200.0)
		assert.NotNil(t, err)
		assert.Equal(t, entity.ErrInsufficientBalance, err.Error())
	})
}
