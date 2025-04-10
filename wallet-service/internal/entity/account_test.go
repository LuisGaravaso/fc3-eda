package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewAccount(t *testing.T) {
	client, _ := NewClient("John", "j@email.com")
	account, err := NewAccount(client)
	assert.NoError(t, err)
	assert.NotEmpty(t, account.Id)
	assert.Equal(t, client, account.Client)
	assert.Equal(t, 0.0, account.Balance)
	assert.NotEmpty(t, account.CreatedAt)
	assert.NotEmpty(t, account.UpdatedAt)
}

func TestCreateNewAccount_MustFailWhenClientIsNil(t *testing.T) {
	account, err := NewAccount(nil)
	assert.NotNil(t, err)
	assert.Nil(t, account)
	assert.Equal(t, ErrInvalidClient, err.Error())
}

func TestAccountCredit(t *testing.T) {
	client, _ := NewClient("John", "j@email.com")
	account, _ := NewAccount(client)
	account.Credit(100.0)
	assert.Equal(t, 100.0, account.Balance)
}

func TestAccountDebit(t *testing.T) {
	client, _ := NewClient("John", "j@email.com")
	account, _ := NewAccount(client)
	account.Credit(100.0)
	account.Debit(50.0)
	assert.Equal(t, 50.0, account.Balance)
}

func TestAccountDebit_MustFailWhenBalanceIsInsufficient(t *testing.T) {
	client, _ := NewClient("John", "j@email.com")
	account, _ := NewAccount(client)
	account.Credit(100.0)
	err := account.Debit(150.0)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInsufficientBalance, err.Error())
}
