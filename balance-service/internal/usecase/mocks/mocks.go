package mocks

import (
	"balance/internal/entity"

	"github.com/stretchr/testify/mock"
)

type BalanceGatewayMock struct {
	mock.Mock
}

func (m *BalanceGatewayMock) FindById(id string) (*entity.AccountBalance, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AccountBalance), args.Error(1)
}

func (m *BalanceGatewayMock) Save(account *entity.AccountBalance) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *BalanceGatewayMock) UpdateBalance(account *entity.AccountBalance) error {
	args := m.Called(account)
	return args.Error(0)
}
