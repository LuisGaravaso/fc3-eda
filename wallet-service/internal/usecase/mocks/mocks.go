package mocks

import (
	"context"
	"time"
	"wallet/internal/entity"
	"wallet/pkg/events"
	"wallet/pkg/uow"

	"github.com/stretchr/testify/mock"
)

// Gateway Mocks

type ClientGateway struct {
	mock.Mock
}

func (m *ClientGateway) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGateway) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

type AccountGateway struct {
	mock.Mock
}

func (m *AccountGateway) FindById(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGateway) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGateway) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

type TransactionGateway struct {
	mock.Mock
}

func (m *TransactionGateway) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

// UOW Mock

type UowMock struct {
	mock.Mock
}

func (m *UowMock) Register(name string, fc uow.RepositoryFactory) {
	m.Called(name, fc)
}

func (m *UowMock) GetRepository(ctx context.Context, name string) (interface{}, error) {
	args := m.Called(ctx, name)
	return args.Get(0), args.Error(1)
}

func (m *UowMock) Do(ctx context.Context, fn func(uow *uow.Uow) error) error {
	args := m.Called(ctx, fn)
	// Execute the function with nil to simulate UOW behavior
	if args.Get(0) == nil {
		_ = fn(nil)
	}
	return args.Error(0)
}

func (m *UowMock) CommitOrRollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) UnRegister(name string) {
	m.Called(name)
}

// Event Mocks

type EventDispatcher struct {
	mock.Mock
}

func (m *EventDispatcher) Register(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *EventDispatcher) Dispatch(event events.EventInterface) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *EventDispatcher) Remove(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *EventDispatcher) Has(eventName string, handler events.EventHandlerInterface) bool {
	args := m.Called(eventName, handler)
	return args.Bool(0)
}

func (m *EventDispatcher) Clear() {
	m.Called()
}

type Event struct {
	mock.Mock
}

func (m *Event) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *Event) GetDateTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *Event) GetPayload() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *Event) SetPayload(payload interface{}) {
	m.Called(payload)
}
