package database

import (
	"balance/internal/entity"
	"database/sql"
	"errors"
)

type BalanceDB struct {
	DB *sql.DB
}

func NewBalanceDB(db *sql.DB) *BalanceDB {
	return &BalanceDB{
		DB: db,
	}
}

func (b *BalanceDB) FindById(id string) (*entity.AccountBalance, error) {
	var balance entity.AccountBalance
	stmt, err := b.DB.Prepare("SELECT account_id, balance FROM account_balances WHERE account_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(&balance.AccountId, &balance.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &balance, nil
}

func (b *BalanceDB) Save(account *entity.AccountBalance) error {
	stmt, err := b.DB.Prepare("INSERT INTO account_balances (account_id, balance) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.AccountId, account.Balance)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceDB) UpdateBalance(account *entity.AccountBalance) error {
	stmt, err := b.DB.Prepare("UPDATE account_balances SET balance = ? WHERE account_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(account.Balance, account.AccountId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}
