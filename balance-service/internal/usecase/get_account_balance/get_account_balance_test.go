package get_account_balance_test

import (
	"balance/internal/entity"
	"balance/internal/usecase/get_account_balance"
	"balance/internal/usecase/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountBalanceUseCase_Execute(t *testing.T) {
	t.Run("should get an existing account balance", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		existingBalance, _ := entity.NewBalance("account1", 100.0)

		balanceMock.On("FindById", "account1").Return(existingBalance, nil)

		useCase := get_account_balance.NewGetAccountBalanceUseCase(balanceMock)

		input := get_account_balance.GetAccountBalanceInputDTO{
			AccountID: "account1",
		}

		output, err := useCase.Execute(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "account1", output.AccountID)
		assert.Equal(t, 100.0, output.Balance)
		balanceMock.AssertExpectations(t)
	})

	t.Run("should return error when account balance not found", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "account1").Return(nil, nil)

		useCase := get_account_balance.NewGetAccountBalanceUseCase(balanceMock)

		input := get_account_balance.GetAccountBalanceInputDTO{
			AccountID: "account1",
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "account balance not found", err.Error())
		balanceMock.AssertExpectations(t)
	})

	t.Run("should return error when gateway returns error", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "account1").Return(nil, errors.New("database error"))

		useCase := get_account_balance.NewGetAccountBalanceUseCase(balanceMock)

		input := get_account_balance.GetAccountBalanceInputDTO{
			AccountID: "account1",
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "database error", err.Error())
		balanceMock.AssertExpectations(t)
	})
}
