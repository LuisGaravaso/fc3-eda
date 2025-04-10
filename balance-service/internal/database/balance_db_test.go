package database_test

import (
	"balance/internal/database"
	"balance/internal/entity"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type BalanceDBTestSuite struct {
	suite.Suite
	DB        *sql.DB
	balanceDB *database.BalanceDB
}

func (s *BalanceDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite", ":memory:")
	s.Nil(err)
	s.DB = db
	s.balanceDB = database.NewBalanceDB(db)
	s.createTable()
}

func (s *BalanceDBTestSuite) createTable() {
	table := `CREATE TABLE account_balances (
        account_id TEXT PRIMARY KEY,
        balance REAL NOT NULL
    );`
	_, err := s.DB.Exec(table)
	s.Nil(err)
}

func (s *BalanceDBTestSuite) TearDownSuite() {
	s.DB.Close()
}

func (s *BalanceDBTestSuite) TearDownTest() {
	s.DB.Exec("DELETE FROM account_balances")
}

func TestBalanceDBTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceDBTestSuite))
}

func (s *BalanceDBTestSuite) TestSave() {
	account, _ := entity.NewBalance("account1", 100.0)
	err := s.balanceDB.Save(account)
	s.Nil(err)

	var accountId string
	var balance float64
	err = s.DB.QueryRow("SELECT account_id, balance FROM account_balances WHERE account_id = ?", account.AccountId).
		Scan(&accountId, &balance)
	s.Nil(err)
	s.Equal("account1", accountId)
	s.Equal(100.0, balance)
}

func (s *BalanceDBTestSuite) TestSaveUpdate_MustFailForDuplicate() {
	account, _ := entity.NewBalance("account1", 100.0)
	err := s.balanceDB.Save(account)
	s.Nil(err)

	account.Balance = 200.0
	err = s.balanceDB.Save(account)

	s.NotNil(err)

}

func (s *BalanceDBTestSuite) TestFindByIdReturnsAccountWhenExists() {
	s.DB.Exec("INSERT INTO account_balances (account_id, balance) VALUES (?, ?)", "account1", 100.0)

	account, err := s.balanceDB.FindById("account1")
	s.Nil(err)
	s.NotNil(account)
	s.Equal("account1", account.AccountId)
	s.Equal(100.0, account.Balance)
}

func (s *BalanceDBTestSuite) TestFindByIdReturnsNilWhenNotExists() {
	account, err := s.balanceDB.FindById("account_not_exists")
	s.Nil(err)
	s.Nil(account)
}

func (s *BalanceDBTestSuite) TestUpdateBalance() {
	s.DB.Exec("INSERT INTO account_balances (account_id, balance) VALUES (?, ?)", "account1", 100.0)

	account, _ := entity.NewBalance("account1", 200.0)
	err := s.balanceDB.UpdateBalance(account)
	s.Nil(err)

	var balance float64
	err = s.DB.QueryRow("SELECT balance FROM account_balances WHERE account_id = ?", account.AccountId).
		Scan(&balance)
	s.Nil(err)
	s.Equal(200.0, balance)
}

func (s *BalanceDBTestSuite) TestUpdateBalanceWithNonExistingAccount() {
	account, _ := entity.NewBalance("account_not_exists", 200.0)
	err := s.balanceDB.UpdateBalance(account)
	s.NotNil(err)
	s.Equal("account not found", err.Error())
}
