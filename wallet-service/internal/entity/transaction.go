package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	ErrInvalidTransaction = "invalid transaction"
	ErrInvalidAccount     = "invalid account"
	ErrInvalidAmount      = "invalid amount"
	ErrNotEnoughBalance   = "not enough balance"
)

type Transaction struct {
	Id          string    `json:"id"`
	AccountFrom *Account  `json:"account_from"`
	AccountTo   *Account  `json:"account_to"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewTransaction(accountFrom, accountTo *Account, amount float64) (*Transaction, error) {
	transaction := &Transaction{
		Id:          uuid.New().String(),
		AccountFrom: accountFrom,
		AccountTo:   accountTo,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}

	err := transaction.Validate()
	if err != nil {
		return nil, err
	}

	transaction.Commit()

	return transaction, nil
}

func (transaction *Transaction) Validate() error {
	if transaction.AccountFrom == nil || transaction.AccountTo == nil {
		return errors.New(ErrInvalidAccount)
	}
	if transaction.Amount <= 0 {
		return errors.New(ErrInvalidAmount)
	}
	if transaction.AccountFrom == transaction.AccountTo {
		return errors.New(ErrInvalidTransaction)
	}
	if transaction.AccountFrom.Balance < transaction.Amount {
		return errors.New(ErrNotEnoughBalance)
	}
	return nil
}

func (transaction *Transaction) Commit() {

	transaction.AccountFrom.Debit(transaction.Amount)
	transaction.AccountTo.Credit(transaction.Amount)
}
