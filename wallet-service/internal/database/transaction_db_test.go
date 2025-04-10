package database

import (
	"database/sql"
	"testing"
	"wallet/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDB *TransactionDB
	accountDB     *AccountDB
	clientDB      *ClientDB
	client1       *entity.Client
	client2       *entity.Client
	account1      *entity.Account
	account2      *entity.Account
}

func (suite *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db = db

	// Create necessary tables for testing
	db.Exec(`CREATE TABLE clients (
        id varchar(255) PRIMARY KEY, 
        name varchar(255), 
        email varchar(255), 
        created_at date
    )`)

	db.Exec(`CREATE TABLE accounts (
        id varchar(255) PRIMARY KEY, 
        client_id varchar(255), 
        balance float, 
        created_at date,
        FOREIGN KEY (client_id) REFERENCES clients(id)
    )`)

	db.Exec(`CREATE TABLE transactions (
        id varchar(255) PRIMARY KEY,
        account_id_from varchar(255),
        account_id_to varchar(255),
        amount float,
        created_at date,
        FOREIGN KEY (account_id_from) REFERENCES accounts(id),
        FOREIGN KEY (account_id_to) REFERENCES accounts(id)
    )`)

	suite.transactionDB = NewTransactionDB(suite.db)
	suite.accountDB = NewAccountDB(suite.db)
	suite.clientDB = NewClientDB(suite.db)
}

func (suite *TransactionDBTestSuite) SetupTest() {
	// Clean tables before each test
	suite.db.Exec("DELETE FROM transactions")
	suite.db.Exec("DELETE FROM accounts")
	suite.db.Exec("DELETE FROM clients")

	// Create and save test clients
	suite.client1, _ = entity.NewClient("John", "john@example.com")
	suite.client2, _ = entity.NewClient("Jane", "jane@example.com")
	suite.clientDB.Save(suite.client1)
	suite.clientDB.Save(suite.client2)

	// Create and save test accounts
	suite.account1, _ = entity.NewAccount(suite.client1)
	suite.account2, _ = entity.NewAccount(suite.client2)
	suite.account1.Credit(1000) // Add initial balance
	suite.accountDB.Save(suite.account1)
	suite.accountDB.Save(suite.account2)
}

func (suite *TransactionDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE transactions")
	suite.db.Exec("DROP TABLE accounts")
	suite.db.Exec("DROP TABLE clients")
}

func (suite *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(suite.account1, suite.account2, 100)
	assert.Nil(suite.T(), err)

	err = suite.transactionDB.Create(transaction)
	assert.Nil(suite.T(), err)

	// Verify the transaction was saved by querying the database directly
	var count int
	row := suite.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", transaction.Id)
	err = row.Scan(&count)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Check if transaction details were correctly saved
	var id, accountFromId, accountToId string
	var amount float64
	row = suite.db.QueryRow(
		"SELECT id, account_id_from, account_id_to, amount FROM transactions WHERE id = ?",
		transaction.Id,
	)
	err = row.Scan(&id, &accountFromId, &accountToId, &amount)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), transaction.Id, id)
	assert.Equal(suite.T(), suite.account1.Id, accountFromId)
	assert.Equal(suite.T(), suite.account2.Id, accountToId)
	assert.Equal(suite.T(), 100.0, amount)

	// Verify account balances were updated
	assert.Equal(suite.T(), 900.0, suite.account1.Balance)
	assert.Equal(suite.T(), 100.0, suite.account2.Balance)
}

func (suite *TransactionDBTestSuite) TestCreateWithInvalidAccounts() {
	// Create transaction with non-existent accounts
	nonExistentAccount, _ := entity.NewAccount(suite.client1)

	// This transaction should save to DB (though in real app this would be inside a transaction)
	transaction, _ := entity.NewTransaction(suite.account1, nonExistentAccount, 100)
	err := suite.transactionDB.Create(transaction)

	// We should get a foreign key constraint error
	assert.NotNil(suite.T(), err)

	// Verify no transaction was inserted
	var count int
	row := suite.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", transaction.Id)
	err = row.Scan(&count)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *TransactionDBTestSuite) TestCreateDuplicate() {
	// Create and save a transaction
	transaction, _ := entity.NewTransaction(suite.account1, suite.account2, 100)
	err := suite.transactionDB.Create(transaction)
	assert.Nil(suite.T(), err)

	// Recreate the balances for testing the duplicate insert
	suite.account1.Credit(100) // Reset balance

	// Try to save the same transaction again
	err = suite.transactionDB.Create(transaction)
	assert.NotNil(suite.T(), err) // Should fail due to primary key constraint

	// Verify only one transaction was inserted
	var count int
	row := suite.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", transaction.Id)
	err = row.Scan(&count)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}
