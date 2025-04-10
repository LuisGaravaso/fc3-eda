package createaccount

import (
	"errors"
	"testing"
	"wallet/internal/entity"
	"wallet/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	mockClient, _ := entity.NewClient("John", "john@example.com")

	mockClientGateway := &mocks.ClientGateway{}
	mockClientGateway.On("Get", mock.Anything).Return(mockClient, nil)

	mockAccountGateway := &mocks.AccountGateway{}
	mockAccountGateway.On("Save", mock.Anything).Return(nil)

	useCase := NewCreateAccountUseCase(mockAccountGateway, mockClientGateway)

	input := CreateAccountInputDTO{
		ClientId: "any_client_id",
	}

	output, err := useCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.Id)
	mockClientGateway.AssertExpectations(t)
	mockAccountGateway.AssertExpectations(t)
	mockClientGateway.AssertNumberOfCalls(t, "Get", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateAccountUseCase_ExecuteWithClientGatewayError(t *testing.T) {
	mockClientGateway := &mocks.ClientGateway{}
	mockClientGateway.On("Get", mock.Anything).Return((*entity.Client)(nil), errors.New("client not found"))

	mockAccountGateway := &mocks.AccountGateway{}

	useCase := NewCreateAccountUseCase(mockAccountGateway, mockClientGateway)

	input := CreateAccountInputDTO{
		ClientId: "any_client_id",
	}

	output, err := useCase.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "client not found", err.Error())
	mockClientGateway.AssertExpectations(t)
	mockClientGateway.AssertNumberOfCalls(t, "Get", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateAccountUseCase_ExecuteWithInvalidClient(t *testing.T) {
	// Setting up a nil client to force NewAccount validation error
	var mockClient *entity.Client = nil

	mockClientGateway := &mocks.ClientGateway{}
	mockClientGateway.On("Get", mock.Anything).Return(mockClient, nil)

	mockAccountGateway := &mocks.AccountGateway{}

	useCase := NewCreateAccountUseCase(mockAccountGateway, mockClientGateway)

	input := CreateAccountInputDTO{
		ClientId: "any_client_id",
	}

	output, err := useCase.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInvalidClient, err.Error())
	mockClientGateway.AssertExpectations(t)
	mockClientGateway.AssertNumberOfCalls(t, "Get", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateAccountUseCase_ExecuteWithSaveError(t *testing.T) {
	mockClient, _ := entity.NewClient("John", "john@example.com")

	mockClientGateway := &mocks.ClientGateway{}
	mockClientGateway.On("Get", mock.Anything).Return(mockClient, nil)

	mockAccountGateway := &mocks.AccountGateway{}
	mockAccountGateway.On("Save", mock.Anything).Return(errors.New("error saving account"))

	useCase := NewCreateAccountUseCase(mockAccountGateway, mockClientGateway)

	input := CreateAccountInputDTO{
		ClientId: "any_client_id",
	}

	output, err := useCase.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "error saving account", err.Error())
	mockClientGateway.AssertExpectations(t)
	mockAccountGateway.AssertExpectations(t)
	mockClientGateway.AssertNumberOfCalls(t, "Get", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Save", 1)
}
