package update_account_balance_test

import (
	"balance/internal/entity"
	"balance/internal/usecase/mocks"
	"balance/internal/usecase/update_account_balance"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateAccountBalanceUseCase_Execute(t *testing.T) {
	t.Run("should update an existing account balance", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		existingBalance, _ := entity.NewBalance("account1", 100.0)

		balanceMock.On("FindById", "account1").Return(existingBalance, nil)
		balanceMock.On("UpdateBalance", mock.MatchedBy(func(acc *entity.AccountBalance) bool {
			return acc.AccountId == "account1" && acc.Balance == 200.0
		})).Return(nil)

		useCase := update_account_balance.NewUpdateAccountBalanceUseCase(balanceMock)

		input := update_account_balance.UpdateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   200.0,
		}

		output, err := useCase.Execute(input)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "account1", output.AccountID)
		assert.Equal(t, 200.0, output.Balance)
		balanceMock.AssertExpectations(t)
	})

	t.Run("should return error when account balance not found", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "account1").Return(nil, errors.New("not found"))

		useCase := update_account_balance.NewUpdateAccountBalanceUseCase(balanceMock)

		input := update_account_balance.UpdateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   200.0,
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "not found", err.Error())
		balanceMock.AssertNotCalled(t, "UpdateBalance")
	})

	t.Run("should return error when account balance is nil", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		balanceMock.On("FindById", "account1").Return(nil, nil)

		useCase := update_account_balance.NewUpdateAccountBalanceUseCase(balanceMock)

		input := update_account_balance.UpdateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   200.0,
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "account balance not found", err.Error())
		balanceMock.AssertNotCalled(t, "UpdateBalance")
	})

	t.Run("should return error when balance is negative", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		existingBalance, _ := entity.NewBalance("account1", 100.0)

		balanceMock.On("FindById", "account1").Return(existingBalance, nil)

		useCase := update_account_balance.NewUpdateAccountBalanceUseCase(balanceMock)

		input := update_account_balance.UpdateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   -50.0,
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, entity.ErrInsufficientBalance, err.Error())
		balanceMock.AssertNotCalled(t, "UpdateBalance")
	})

	t.Run("should return error when update balance operation fails", func(t *testing.T) {
		balanceMock := &mocks.BalanceGatewayMock{}
		existingBalance, _ := entity.NewBalance("account1", 100.0)

		balanceMock.On("FindById", "account1").Return(existingBalance, nil)
		balanceMock.On("UpdateBalance", mock.Anything).Return(errors.New("database error"))

		useCase := update_account_balance.NewUpdateAccountBalanceUseCase(balanceMock)

		input := update_account_balance.UpdateAccountBalanceInputDTO{
			AccountID: "account1",
			Balance:   200.0,
		}

		output, err := useCase.Execute(input)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, "database error", err.Error())
	})
}
