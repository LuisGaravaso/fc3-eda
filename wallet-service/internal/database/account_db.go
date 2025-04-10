package database

import (
	"database/sql"
	"fmt"
	"wallet/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{DB: db}
}

func (a *AccountDB) FindById(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	query := `SELECT 
				a.id, 
				a.client_id, 
				a.balance, 
				a.created_at, 
				c.id, 
				c.name, 
				c.email, 
				c.created_at 
			  FROM accounts a INNER JOIN clients c 
			  ON a.client_id = c.id 
			  WHERE a.id = ?`

	row := a.DB.QueryRow(query, id)
	err := row.Scan(
		&account.Id,
		&account.Client.Id,
		&account.Balance,
		&account.CreatedAt,
		&client.Id,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil

}

func (a *AccountDB) Save(account *entity.Account) error {
	// Check if the client exists
	clientQuery := `SELECT id FROM clients WHERE id = ?`
	var clientId string
	err := a.DB.QueryRow(clientQuery, account.Client.Id).Scan(&clientId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("client does not exist")
	}

	// Check if the account already exists
	accountQuery := `SELECT id FROM accounts WHERE id = ?`
	var accountId string
	err = a.DB.QueryRow(accountQuery, account.Id).Scan(&accountId)
	if err == nil {
		return fmt.Errorf("account already exists")
	}
	// Insert the account
	insertQuery := `INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)`
	a.DB.Exec(insertQuery, account.Id, account.Client.Id, account.Balance, account.CreatedAt)
	return nil
}

func (a *AccountDB) UpdateBalance(account *entity.Account) error {
	updateQuery := `UPDATE accounts SET balance = ? WHERE id = ?`
	_, err := a.DB.Exec(updateQuery, account.Balance, account.Id)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}
	return nil
}
