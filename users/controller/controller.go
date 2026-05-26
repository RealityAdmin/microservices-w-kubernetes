package controller

import (
	"context"
	"errors"

	"k-microserv-kuber.com/users"
	"k-microserv-kuber.com/users/database"
)

type Repo interface {
	GetUser(ctx context.Context, id int) (*users.User, error)
	// PutUser(ctx context.Context, id int, user *users.User) error
}

type Controller struct {
	repo Repo
}

func New(repo Repo) *Controller {
	return &Controller{repo: repo}
}

var ErrNotFound = errors.New("User not found")

func (c *Controller) GetUser(ctx context.Context, id int) (*users.User, error) {
	res, err := c.repo.GetUser(ctx, id)
	if err != nil && errors.Is(err, database.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, err
}
