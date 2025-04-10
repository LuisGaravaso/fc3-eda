package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	ErrInvalidName     = "invalid name"
	ErrInvalidEmail    = "invalid email"
	ErrAccountMismatch = "account client mismatch"
)

type Client struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Accounts  []*Account `json:"accounts"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewClient(name, email string) (*Client, error) {
	client := &Client{
		Id:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := client.Validate()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Validate() error {
	if c.Name == "" {
		return errors.New(ErrInvalidName)
	}
	if c.Email == "" {
		return errors.New(ErrInvalidEmail)
	}
	return nil
}

func (c *Client) Update(name, email string) error {
	c.Name = name
	c.Email = email
	c.UpdatedAt = time.Now()

	err := c.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddAccount(account *Account) error {
	if account.Client.Id != c.Id {
		return errors.New(ErrAccountMismatch)
	}
	c.Accounts = append(c.Accounts, account)
	return nil
}
