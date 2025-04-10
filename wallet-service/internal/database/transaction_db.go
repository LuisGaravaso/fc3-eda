package database

import (
	"database/sql"
	"fmt"
	"wallet/internal/entity"
)

type TransactionDB struct {
	DB *sql.DB
}

func NewTransactionDB(db *sql.DB) *TransactionDB {
	return &TransactionDB{DB: db}
}

func (t *TransactionDB) Create(transaction *entity.Transaction) error {
	// Check if the account_from exists
	var accountFromExists bool
	t.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM accounts WHERE id = ?)`, transaction.AccountFrom.Id).Scan(&accountFromExists)
	if !accountFromExists {
		return fmt.Errorf("account_from with id %s does not exist", transaction.AccountFrom.Id)
	}

	// Check if the account_to exists
	var accountToExists bool
	t.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM accounts WHERE id = ?)`, transaction.AccountTo.Id).Scan(&accountToExists)

	if !accountToExists {
		return fmt.Errorf("account_to with id %s does not exist", transaction.AccountTo.Id)
	}

	// Proceed with the transaction creation
	query := `INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := t.DB.Exec(query,
		transaction.Id,
		transaction.AccountFrom.Id,
		transaction.AccountTo.Id,
		transaction.Amount,
		transaction.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
