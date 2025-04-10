package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("John", "john@email.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, client.Id)
	assert.Equal(t, "John", client.Name)
	assert.Equal(t, "john@email.com", client.Email)
}

func TestCreateNewClient_MustFailWhenParamsAreInvalid(t *testing.T) {
	client, err := NewClient("", "john@email.com")
	assert.NotNil(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidName, err.Error())
}

func TestCreateNewClient_MustFailWhenEmailIsInvalid(t *testing.T) {
	client, err := NewClient("John", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidEmail, err.Error())
}

func TestUpdateClient(t *testing.T) {
	c, _ := NewClient("John", "john@email.com")
	err := c.Update("John Updated", "john@email.com")
	assert.NoError(t, err)
	assert.Equal(t, "John Updated", c.Name)
	assert.Equal(t, "john@email.com", c.Email)
}

func TestUpdateClient_MustFailWhenParamsAreInvalid(t *testing.T) {
	c, _ := NewClient("John", "john@email.com")
	err := c.Update("", "john@email.com")
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidName, err.Error())
}

func TestAddAccount(t *testing.T) {
	c, _ := NewClient("John", "john@email.com")
	account, _ := NewAccount(c)
	c.AddAccount(account)
	assert.Equal(t, 1, len(c.Accounts))
	assert.Equal(t, account, c.Accounts[0])
}

func TestAddAccount_MustFailWhenAccountBelongsToOtherClient(t *testing.T) {
	c1, _ := NewClient("John", "john@email.com")
	a1, _ := NewAccount(c1)

	c2, _ := NewClient("Jane", "jane@email.com")
	err := c2.AddAccount(a1)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAccountMismatch, err.Error())
	assert.Equal(t, 0, len(c2.Accounts))
}
