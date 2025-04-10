package createaccount

import (
	"wallet/internal/entity"
	"wallet/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientId string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	Id string `json:"id"`
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(accountGateway gateway.AccountGateway, clientGateway gateway.ClientGateway) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: accountGateway,
		ClientGateway:  clientGateway,
	}
}

func (uc *CreateAccountUseCase) Execute(input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	client, err := uc.ClientGateway.Get(input.ClientId)
	if err != nil {
		return nil, err
	}

	account, err := entity.NewAccount(client)
	if err != nil {
		return nil, err
	}

	err = uc.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}

	output := &CreateAccountOutputDTO{
		Id: account.Id,
	}

	return output, nil
}
