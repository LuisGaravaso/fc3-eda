package update_account_balance

import (
	"balance/internal/gateway"
	"errors"
)

type UpdateAccountBalanceInputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type UpdateAccountBalanceOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type UpdateAccountBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewUpdateAccountBalanceUseCase(balanceGateway gateway.BalanceGateway) *UpdateAccountBalanceUseCase {
	return &UpdateAccountBalanceUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *UpdateAccountBalanceUseCase) Execute(input UpdateAccountBalanceInputDTO) (*UpdateAccountBalanceOutputDTO, error) {
	// Find the existing account balance
	existingBalance, err := uc.BalanceGateway.FindById(input.AccountID)
	if err != nil {
		return nil, err
	}
	if existingBalance == nil {
		return nil, errors.New("account balance not found")
	}

	// Update the balance
	err = existingBalance.UpdateBalance(input.Balance)
	if err != nil {
		return nil, err
	}

	// Save the updated account balance
	err = uc.BalanceGateway.UpdateBalance(existingBalance)
	if err != nil {
		return nil, err
	}

	// Return the updated account balance
	return &UpdateAccountBalanceOutputDTO{
		AccountID: existingBalance.AccountId,
		Balance:   existingBalance.Balance,
	}, nil
}
