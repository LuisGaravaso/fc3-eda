package createclient

import (
	"time"
	"wallet/internal/entity"
	"wallet/internal/gateway"
)

type CreateClientInputDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateClientOutputDTO struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateClientUseCase struct {
	ClientGateway gateway.ClientGateway
}

func NewCreateClientUseCase(clientGateway gateway.ClientGateway) *CreateClientUseCase {
	return &CreateClientUseCase{
		ClientGateway: clientGateway,
	}
}

func (uc *CreateClientUseCase) Execute(input CreateClientInputDTO) (*CreateClientOutputDTO, error) {
	client, err := entity.NewClient(input.Name, input.Email)
	if err != nil {
		return nil, err
	}

	err = uc.ClientGateway.Save(client)
	if err != nil {
		return nil, err
	}

	output := &CreateClientOutputDTO{
		Id:        client.Id,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}

	return output, nil
}
