package gateway

import "wallet/internal/entity"

type AccountGateway interface {
	FindById(id string) (*entity.Account, error)
	Save(account *entity.Account) error
	UpdateBalance(account *entity.Account) error
}
