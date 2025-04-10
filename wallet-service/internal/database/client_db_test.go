package database

import (
	"database/sql"
	"testing"
	"wallet/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (suite *ClientDBTestSuite) SetupSuite() {
	// Initialize the database connection and ClientDB instance here
	// For example, you can use a mock database or an in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db = db
	db.Exec("CREATE TABLE clients (id varchar(255) PRIMARY KEY, name varchar(255), email varchar(255), created_at date)")

	suite.clientDB = NewClientDB(suite.db)
}

func (suite *ClientDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE clients")
}

func (suite *ClientDBTestSuite) TestSave() {
	client, _ := entity.NewClient("John Doe", "john@example.com")
	err := suite.clientDB.Save(client)

	assert.Nil(suite.T(), err)

	// Verify the client was saved by querying the database directly
	var count int
	row := suite.db.QueryRow("SELECT COUNT(*) FROM clients WHERE id = ?", client.Id)
	err = row.Scan(&count)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

func (suite *ClientDBTestSuite) TestSaveDuplicate() {
	client, _ := entity.NewClient("Jane Doe", "jane@example.com")
	err := suite.clientDB.Save(client)
	assert.Nil(suite.T(), err)

	// Try to save the same client again
	err = suite.clientDB.Save(client)
	assert.NotNil(suite.T(), err) // Should fail due to primary key constraint
}

func (suite *ClientDBTestSuite) TestGet() {
	// Create and save a client first
	expectedClient, _ := entity.NewClient("Alice Smith", "alice@example.com")
	suite.db.Exec(
		"INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		expectedClient.Id, expectedClient.Name, expectedClient.Email, expectedClient.CreatedAt,
	)

	// Retrieve the client
	client, err := suite.clientDB.Get(expectedClient.Id)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), expectedClient.Id, client.Id)
	assert.Equal(suite.T(), expectedClient.Name, client.Name)
	assert.Equal(suite.T(), expectedClient.Email, client.Email)
	// Note: Time comparison might need some flexibility due to serialization differences
}

func (suite *ClientDBTestSuite) TestGetNonExistentClient() {
	client, err := suite.clientDB.Get("non-existent-id")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrNoRows, err)
	assert.Nil(suite.T(), client)
}

func (suite *ClientDBTestSuite) SetupTest() {
	// Clean up the table before each test
	suite.db.Exec("DELETE FROM clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}
