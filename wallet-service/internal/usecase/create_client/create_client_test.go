package createclient

import (
	"errors"
	"testing"
	"wallet/internal/entity"
	"wallet/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClientUseCase_Execute(t *testing.T) {
	m := &mocks.ClientGateway{}
	m.On("Save", mock.Anything).Return(nil)

	useCase := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	output, err := useCase.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.Name, output.Name)
	assert.Equal(t, input.Email, output.Email)
	assert.NotEmpty(t, output.Id)
	assert.NotEmpty(t, output.CreatedAt)
	assert.NotEmpty(t, output.UpdatedAt)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateClientUseCase_ExecuteWithInvalidName(t *testing.T) {
	m := &mocks.ClientGateway{}

	useCase := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "",
		Email: "john@example.com",
	}

	output, err := useCase.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInvalidName, err.Error())
	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateClientUseCase_ExecuteWithInvalidEmail(t *testing.T) {
	m := &mocks.ClientGateway{}

	useCase := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "John Doe",
		Email: "",
	}

	output, err := useCase.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInvalidEmail, err.Error())
	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateClientUseCase_ExecuteWithSaveError(t *testing.T) {
	m := &mocks.ClientGateway{}
	m.On("Save", mock.Anything).Return(errors.New("error saving client"))

	useCase := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	output, err := useCase.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "error saving client", err.Error())
	m.AssertNumberOfCalls(t, "Save", 1)
	m.AssertExpectations(t)
}
