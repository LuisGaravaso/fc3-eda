package entity

import "errors"

type AccountBalance struct {
	AccountId string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

const (
	ErrInvalidClient       = "invalid client"
	ErrInsufficientBalance = "insufficient balance"
)

func NewBalance(accountId string, balance float64) (*AccountBalance, error) {
	accBalance := &AccountBalance{
		AccountId: accountId,
		Balance:   balance,
	}

	if err := accBalance.Validate(); err != nil {
		return nil, err
	}

	return accBalance, nil

}

func (b *AccountBalance) Validate() error {
	if b.AccountId == "" {
		return errors.New(ErrInvalidClient)
	}
	if b.Balance < 0 {
		return errors.New(ErrInsufficientBalance)
	}
	return nil
}

func (b *AccountBalance) UpdateBalance(newBalance float64) error {
	if newBalance < 0 {
		return errors.New(ErrInsufficientBalance)
	}
	b.Balance = newBalance
	return nil
}
