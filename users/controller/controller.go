package controller

import (
	"context"
	"errors"

	"k-microserv-kuber.com/users"
	"k-microserv-kuber.com/users/database"
)

type Database interface {
	GetUser(ctx context.Context, userId int) (*users.User, error)
	PutUser(ctx context.Context, userId int, user *users.User) error
}

type Controller struct {
	db Database
}

func New(db Database) *Controller {
	return &Controller{db: db}
}

var ErrNotFound = errors.New("User not found")

func (c *Controller) GetUser(ctx context.Context, userId int) (*users.User, error) {
	res, err := c.db.GetUser(ctx, userId)
	if err != nil && errors.Is(err, database.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, err
}

func (c *Controller) PutUser(ctx context.Context, userId int, user *users.User) error {
	return c.db.PutUser(ctx, userId, user)
}
