package gateway

import (
	"balance/internal/entity"
)

type BalanceGateway interface {
	FindById(id string) (*entity.AccountBalance, error)
	Save(account *entity.AccountBalance) error
	UpdateBalance(account *entity.AccountBalance) error
}
