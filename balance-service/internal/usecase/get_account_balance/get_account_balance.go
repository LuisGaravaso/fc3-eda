package get_account_balance

import (
	"balance/internal/gateway"
	"errors"
)

type GetAccountBalanceInputDTO struct {
	AccountID string `json:"account_id"`
}

type GetAccountBalanceOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type GetAccountBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewGetAccountBalanceUseCase(balanceGateway gateway.BalanceGateway) *GetAccountBalanceUseCase {
	return &GetAccountBalanceUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *GetAccountBalanceUseCase) Execute(input GetAccountBalanceInputDTO) (*GetAccountBalanceOutputDTO, error) {
	// Find the account balance
	accountBalance, err := uc.BalanceGateway.FindById(input.AccountID)
	if err != nil {
		return nil, err
	}
	if accountBalance == nil {
		return nil, errors.New("account balance not found")
	}

	// Return the account balance
	return &GetAccountBalanceOutputDTO{
		AccountID: accountBalance.AccountId,
		Balance:   accountBalance.Balance,
	}, nil
}
