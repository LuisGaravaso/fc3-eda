package database

import (
	"database/sql"
	"wallet/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{DB: db}
}

func (c *ClientDB) Get(id string) (*entity.Client, error) {
	query := `SELECT id, name, email, created_at FROM clients WHERE id = ?`
	row := c.DB.QueryRow(query, id)

	client := &entity.Client{}
	err := row.Scan(&client.Id, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientDB) Save(client *entity.Client) error {
	query := `INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)`
	_, err := c.DB.Exec(query, client.Id, client.Name, client.Email, client.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
