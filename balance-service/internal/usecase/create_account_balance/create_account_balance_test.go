package create_account_balance_test

import (
	"balance/internal/entity"
	"balance/internal/usecase/create_account_balance"
	"balance/internal/usecase/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountBalanceUseCase_Execute(t *testing.T) {
	t.Run("should create a new account balance when it doesn't exist", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "account1").Return(nil, errors.New("not found"))

		balanceMock.On("Save", mock.MatchedBy(func(acc *entity.AccountBalance) bool {
			return acc.AccountId == "account1" && acc.Balance == 100.0
		})).Return(nil)

		useCase := create_account_balance.NewCreateAccountBalanceUseCase(balanceMock)

		input := create_account_balance.CreateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   100.0,
		}

		output, err := useCase.Execute(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "account1", output.AccountID)
		assert.Equal(t, 100.0, output.Balance)
		balanceMock.AssertExpectations(t)
	})

	t.Run("should return existing account balance when it already exists", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		existingBalance, _ := entity.NewBalance("account1", 150.0)
		balanceMock.On("FindById", "account1").Return(existingBalance, nil)

		useCase := create_account_balance.NewCreateAccountBalanceUseCase(balanceMock)

		input := create_account_balance.CreateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   100.0, // Note: This value is ignored since we return the existing balance
		}

		output, err := useCase.Execute(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "account1", output.AccountID)
		assert.Equal(t, 150.0, output.Balance) // Should be the existing balance
		balanceMock.AssertNotCalled(t, "Save")
	})

	t.Run("should return error when account ID is invalid", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "").Return(nil, errors.New("not found"))

		useCase := create_account_balance.NewCreateAccountBalanceUseCase(balanceMock)

		input := create_account_balance.CreateAccountBalanceInputDTO{
			AccountID: "", // Invalid ID
			Balance:   100.0,
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, entity.ErrInvalidClient, err.Error())
		balanceMock.AssertNotCalled(t, "Save")
	})

	t.Run("should return error when balance gateway save fails", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "account1").Return(nil, errors.New("not found"))

		balanceMock.On("Save", mock.Anything).Return(errors.New("database error"))

		useCase := create_account_balance.NewCreateAccountBalanceUseCase(balanceMock)

		input := create_account_balance.CreateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   100.0,
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "database error", err.Error())
	})
}
