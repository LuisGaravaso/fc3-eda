package create_account_balance

import (
	"balance/internal/entity"
	"balance/internal/gateway"
)

type CreateAccountBalanceInputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type CreateAccountBalanceOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type CreateAccountBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewCreateAccountBalanceUseCase(balanceGateway gateway.BalanceGateway) *CreateAccountBalanceUseCase {
	return &CreateAccountBalanceUseCase{
		BalanceGateway: balanceGateway,
	}
}

func (uc *CreateAccountBalanceUseCase) Execute(input CreateAccountBalanceInputDTO) (*CreateAccountBalanceOutputDTO, error) {
	// Check if account balance already exists
	existingBalance, err := uc.BalanceGateway.FindById(input.AccountID)
	if err == nil && existingBalance != nil {
		// Return existing balance if found
		return &CreateAccountBalanceOutputDTO{
			AccountID: existingBalance.AccountId,
			Balance:   existingBalance.Balance,
		}, nil
	}

	// Create new account balance if it doesn't exist
	accountBalance, err := entity.NewBalance(input.AccountID, input.Balance)
	if err != nil {
		return nil, err
	}

	// Save the new account balance
	err = uc.BalanceGateway.Save(accountBalance)
	if err != nil {
		return nil, err
	}

	// Return the created account balance
	return &CreateAccountBalanceOutputDTO{
		AccountID: accountBalance.AccountId,
		Balance:   accountBalance.Balance,
	}, nil
}
