package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	ErrInvalidClient       = "invalid client"
	ErrInsufficientBalance = "insufficient balance"
)

type Account struct {
	Id        string    `json:"id"`
	Client    *Client   `json:"client"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAccount(client *Client) (*Account, error) {
	account := &Account{
		Id:        uuid.New().String(),
		Client:    client,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := account.Validate()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (account *Account) Validate() error {
	if account.Client == nil {
		return errors.New(ErrInvalidClient)
	}
	return nil
}

func (a *Account) Credit(amount float64) {
	a.Balance += amount
	a.UpdatedAt = time.Now()
}

func (a *Account) Debit(amount float64) error {
	if a.Balance < amount {
		return errors.New(ErrInsufficientBalance)
	}
	a.Balance -= amount
	a.UpdatedAt = time.Now()
	return nil
}
