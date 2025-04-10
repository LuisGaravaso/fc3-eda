package database

import (
	"database/sql"
	"testing"
	"wallet/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
	clientDB  *ClientDB
}

func (suite *AccountDBTestSuite) SetupSuite() {
	// Initialize the database connection and AccountDB instance here
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db = db

	// Create clients table
	db.Exec(`CREATE TABLE clients (
        id varchar(255) PRIMARY KEY, 
        name varchar(255), 
        email varchar(255), 
        created_at date
    )`)

	// Create accounts table
	db.Exec(`CREATE TABLE accounts (
        id varchar(255) PRIMARY KEY, 
        client_id varchar(255), 
        balance float, 
        created_at date,
        FOREIGN KEY (client_id) REFERENCES clients(id)
    )`)

	suite.accountDB = NewAccountDB(suite.db)
	suite.clientDB = NewClientDB(suite.db)
}

func (suite *AccountDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE accounts")
	suite.db.Exec("DROP TABLE clients")
}

func (suite *AccountDBTestSuite) SetupTest() {
	// Clean up the tables before each test
	suite.db.Exec("DELETE FROM accounts")
	suite.db.Exec("DELETE FROM clients")
}

func (suite *AccountDBTestSuite) TestSave() {
	// First, create and save a client
	client, _ := entity.NewClient("John Doe", "john@example.com")
	err := suite.clientDB.Save(client)
	assert.Nil(suite.T(), err)

	// Then create and save an account for that client
	account, _ := entity.NewAccount(client)
	err = suite.accountDB.Save(account)
	assert.Nil(suite.T(), err)

	// Verify the account was saved by querying the database directly
	var count int
	row := suite.db.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = ?", account.Id)
	err = row.Scan(&count)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

func (suite *AccountDBTestSuite) TestSaveDuplicate() {
	// First, create and save a client
	client, _ := entity.NewClient("Jane Doe", "jane@example.com")
	err := suite.clientDB.Save(client)
	assert.Nil(suite.T(), err)

	// Create and save an account
	account, _ := entity.NewAccount(client)
	err = suite.accountDB.Save(account)
	assert.Nil(suite.T(), err)

	// Try to save the same account again
	err = suite.accountDB.Save(account)
	assert.NotNil(suite.T(), err) // Should fail due to primary key constraint
}

func (suite *AccountDBTestSuite) TestFindById() {
	// First, create and save a client
	client, _ := entity.NewClient("Alice Smith", "alice@example.com")
	err := suite.clientDB.Save(client)
	assert.Nil(suite.T(), err)

	// Then create and save an account
	expectedAccount, _ := entity.NewAccount(client)
	expectedAccount.Credit(100.0) // Add some balance
	err = suite.accountDB.Save(expectedAccount)
	assert.Nil(suite.T(), err)

	// Retrieve the account
	account, err := suite.accountDB.FindById(expectedAccount.Id)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), account)
	assert.Equal(suite.T(), expectedAccount.Id, account.Id)
	assert.Equal(suite.T(), expectedAccount.Balance, account.Balance)
	assert.NotNil(suite.T(), account.Client)
	assert.Equal(suite.T(), client.Id, account.Client.Id)
	assert.Equal(suite.T(), client.Name, account.Client.Name)
	assert.Equal(suite.T(), client.Email, account.Client.Email)
}

func (suite *AccountDBTestSuite) TestFindByIdNonExistent() {
	account, err := suite.accountDB.FindById("non-existent-id")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrNoRows, err)
	assert.Nil(suite.T(), account)
}

func (suite *AccountDBTestSuite) TestSaveWithNonExistentClient() {
	// Create a client but don't save it to the database
	client, _ := entity.NewClient("Bob Johnson", "bob@example.com")
	account, _ := entity.NewAccount(client)

	// Try to save the account with a client that doesn't exist in the database
	err := suite.accountDB.Save(account)
	assert.Error(suite.T(), err)
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}
