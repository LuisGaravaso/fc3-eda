package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewTransaction(t *testing.T) {
	client1, _ := NewClient("John", "john@email.com")
	account1, _ := NewAccount(client1)
	account1.Credit(100.0)

	client2, _ := NewClient("Jane", "jane@email.com")
	account2, _ := NewAccount(client2)

	transaction, err := NewTransaction(account1, account2, 50.0)

	assert.NoError(t, err)
	assert.NotEmpty(t, transaction.Id)
	assert.Equal(t, account1, transaction.AccountFrom)
	assert.Equal(t, account2, transaction.AccountTo)
	assert.Equal(t, 50.0, transaction.Amount)
	assert.NotEmpty(t, transaction.CreatedAt)
	assert.Equal(t, 50.0, account1.Balance)
	assert.Equal(t, 50.0, account2.Balance)
}

func TestCreateNewTransaction_MustFailWhenAccountFromIsInvalid(t *testing.T) {
	client2, _ := NewClient("Jane", "jane@email.com")
	account2, _ := NewAccount(client2)

	transaction, err := NewTransaction(nil, account2, 50.0)

	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidAccount, err.Error())
}

func TestCreateNewTransaction_MustFailWhenAccountToIsInvalid(t *testing.T) {
	client1, _ := NewClient("John", "john@email.com")
	account1, _ := NewAccount(client1)
	account1.Credit(100.0)

	transaction, err := NewTransaction(account1, nil, 50.0)

	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidAccount, err.Error())
}

func TestCreateNewTransaction_MustFailWhenAmountIsZeroOrNegative(t *testing.T) {
	client1, _ := NewClient("John", "john@email.com")
	account1, _ := NewAccount(client1)
	account1.Credit(100.0)

	client2, _ := NewClient("Jane", "jane@email.com")
	account2, _ := NewAccount(client2)

	transaction, err := NewTransaction(account1, account2, 0)

	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidAmount, err.Error())

	transaction, err = NewTransaction(account1, account2, -10.0)

	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidAmount, err.Error())
}

func TestCreateNewTransaction_MustFailWhenAccountsAreTheSame(t *testing.T) {
	client1, _ := NewClient("John", "john@email.com")
	account1, _ := NewAccount(client1)
	account1.Credit(100.0)

	transaction, err := NewTransaction(account1, account1, 50.0)

	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidTransaction, err.Error())
}

func TestCreateNewTransaction_MustFailWhenBalanceIsInsufficient(t *testing.T) {
	client1, _ := NewClient("John", "john@email.com")
	account1, _ := NewAccount(client1)
	account1.Credit(30.0)

	client2, _ := NewClient("Jane", "jane@email.com")
	account2, _ := NewAccount(client2)

	transaction, err := NewTransaction(account1, account2, 50.0)

	assert.Nil(t, transaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotEnoughBalance, err.Error())
}

func TestCommitTransaction(t *testing.T) {
	client1, _ := NewClient("John", "john@email.com")
	account1, _ := NewAccount(client1)
	account1.Credit(100.0)

	client2, _ := NewClient("Jane", "jane@email.com")
	account2, _ := NewAccount(client2)

	transaction := &Transaction{
		Id:          "any_id",
		AccountFrom: account1,
		AccountTo:   account2,
		Amount:      50.0,
	}

	transaction.Commit()

	assert.Equal(t, 50.0, account1.Balance)
	assert.Equal(t, 50.0, account2.Balance)
}
