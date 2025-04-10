package createtransaction

import (
	"context"
	"errors"
	"testing"
	"wallet/internal/entity"
	"wallet/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John", "john@example.com")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(100)

	client2, _ := entity.NewClient("Jane", "jane@example.com")
	account2, _ := entity.NewAccount(client2)

	mockAccountGateway := &mocks.AccountGateway{}
	mockAccountGateway.On("FindById", "account1").Return(account1, nil)
	mockAccountGateway.On("FindById", "account2").Return(account2, nil)
	mockAccountGateway.On("UpdateBalance", account1).Return(nil)
	mockAccountGateway.On("UpdateBalance", account2).Return(nil)

	mockTransactionGateway := &mocks.TransactionGateway{}
	mockTransactionGateway.On("Create", mock.Anything).Return(nil)

	mockUow := &mocks.UowMock{}
	mockUow.On("GetRepository", mock.Anything, "AccountRepository").Return(mockAccountGateway, nil)
	mockUow.On("GetRepository", mock.Anything, "TransactionRepository").Return(mockTransactionGateway, nil)
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	mockEventDispatcher := &mocks.EventDispatcher{}
	mockEventDispatcher.On("Dispatch", mock.Anything).Return(nil)

	mockEvent := &mocks.Event{}
	mockEvent.On("SetPayload", mock.Anything).Return()

	mockEvent2 := &mocks.Event{}
	mockEvent2.On("SetPayload", mock.Anything).Return()

	useCase := NewCreateTransactionUseCase(mockUow, mockEventDispatcher, mockEvent, mockEvent2)

	input := CreateTransactionInputDTO{
		AccountIdFrom: "account1",
		AccountIdTo:   "account2",
		Amount:        50,
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockEventDispatcher.AssertExpectations(t)
	mockEvent.AssertCalled(t, "SetPayload", mock.Anything)
	mockEventDispatcher.AssertCalled(t, "Dispatch", mockEvent)
}

func TestCreateTransactionUseCase_FailGetAccountRepository(t *testing.T) {
	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(errors.New("error getting repository"))

	mockEventDispatcher := &mocks.EventDispatcher{}
	mockEvent := &mocks.Event{}
	mockEvent2 := &mocks.Event{}

	useCase := NewCreateTransactionUseCase(mockUow, mockEventDispatcher, mockEvent, mockEvent2)

	input := CreateTransactionInputDTO{
		AccountIdFrom: "account1",
		AccountIdTo:   "account2",
		Amount:        50,
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "error getting repository", err.Error())
	mockUow.AssertExpectations(t)
	mockEventDispatcher.AssertNotCalled(t, "Dispatch", mock.Anything)
	mockEvent.AssertNotCalled(t, "SetPayload", mock.Anything)
}

func TestCreateTransactionUseCase_FailGetTransactionRepository(t *testing.T) {

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(errors.New("error getting transaction repository"))

	mockEventDispatcher := &mocks.EventDispatcher{}
	mockEvent := &mocks.Event{}
	mockEvent2 := &mocks.Event{}

	useCase := NewCreateTransactionUseCase(mockUow, mockEventDispatcher, mockEvent, mockEvent2)

	input := CreateTransactionInputDTO{
		AccountIdFrom: "account1",
		AccountIdTo:   "account2",
		Amount:        50,
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "error getting transaction repository", err.Error())
	mockUow.AssertExpectations(t)
	mockEventDispatcher.AssertNotCalled(t, "Dispatch", mock.Anything)
}

func TestCreateTransactionUseCase_AccountFromNotFound(t *testing.T) {
	mockAccountGateway := &mocks.AccountGateway{}
	mockAccountGateway.On("FindById", "account1").Return((*entity.Account)(nil), errors.New("account not found"))

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(errors.New("account not found"))

	mockEventDispatcher := &mocks.EventDispatcher{}
	mockEvent := &mocks.Event{}
	mockEvent2 := &mocks.Event{}
	useCase := NewCreateTransactionUseCase(mockUow, mockEventDispatcher, mockEvent, mockEvent2)

	input := CreateTransactionInputDTO{
		AccountIdFrom: "account1",
		AccountIdTo:   "account2",
		Amount:        50,
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "account not found", err.Error())
	mockUow.AssertExpectations(t)
	mockEventDispatcher.AssertNotCalled(t, "Dispatch", mock.Anything)
}

func TestCreateTransactionUseCase_AccountToNotFound(t *testing.T) {
	client1, _ := entity.NewClient("John", "john@example.com")
	account1, _ := entity.NewAccount(client1)

	mockAccountGateway := &mocks.AccountGateway{}
	mockAccountGateway.On("FindById", "account1").Return(account1, nil)
	mockAccountGateway.On("FindById", "account2").Return((*entity.Account)(nil), errors.New("account not found"))

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(errors.New("account not found"))

	mockEventDispatcher := &mocks.EventDispatcher{}
	mockEvent := &mocks.Event{}
	mockEvent2 := &mocks.Event{}
	useCase := NewCreateTransactionUseCase(mockUow, mockEventDispatcher, mockEvent, mockEvent2)

	input := CreateTransactionInputDTO{
		AccountIdFrom: "account1",
		AccountIdTo:   "account2",
		Amount:        50,
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "account not found", err.Error())
	mockUow.AssertExpectations(t)
	mockEventDispatcher.AssertNotCalled(t, "Dispatch", mock.Anything)
}

func TestCreateTransactionUseCase_TransactionCreationFailed(t *testing.T) {
	client1, _ := entity.NewClient("John", "john@example.com")
	account1, _ := entity.NewAccount(client1)
	account1.Credit(100)

	client2, _ := entity.NewClient("Jane", "jane@example.com")
	account2, _ := entity.NewAccount(client2)

	mockAccountGateway := &mocks.AccountGateway{}
	mockAccountGateway.On("FindById", "account1").Return(account1, nil)
	mockAccountGateway.On("FindById", "account2").Return(account2, nil)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(errors.New(entity.ErrInvalidTransaction))

	mockEventDispatcher := &mocks.EventDispatcher{}
	mockEvent := &mocks.Event{}
	mockEvent2 := &mocks.Event{}
	useCase := NewCreateTransactionUseCase(mockUow, mockEventDispatcher, mockEvent, mockEvent2)

	input := CreateTransactionInputDTO{
		AccountIdFrom: "account1",
		AccountIdTo:   "account1", // Same account - will fail validation
		Amount:        50,
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInvalidTransaction, err.Error())
	mockUow.AssertExpectations(t)
	mockEventDispatcher.AssertNotCalled(t, "Dispatch", mock.Anything)
}
