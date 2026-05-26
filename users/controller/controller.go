package controller

import (
	"context"

	"k-microserv-kuber.com/users"
)

type Repo interface {
	GetUser(ctx context.Context, id int) (*users.User, error)
	PutUser(ctx context.Context, id int, user *users.User) error
}

type Controller struct {
	repo Repo
}

func New(repo Repo) *Controller {
	return &Controller{repo: repo}
}
