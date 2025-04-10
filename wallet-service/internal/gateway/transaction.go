package gateway

import "wallet/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
