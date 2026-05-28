package controller

import (
	"context"
	"errors"

	"k-microserv-kuber.com/products"
	"k-microserv-kuber.com/products/database"
)

type Database interface {
	GetProduct(_ context.Context, productId int) (*products.Product, error)
	InsertProduct(_ context.Context, productId int, product *products.Product) error
}

type Controller struct {
	db Database
}

func New(db Database) *Controller {
	return &Controller{db: db}
}

func (c *Controller) GetProduct(ctx context.Context, productId int) (*products.Product, error) {
	res, err := c.db.GetProduct(ctx, productId)
	if err != nil && errors.Is(err, database.ErrNotFound) {
		return nil, database.ErrNotFound
	}
	return res, err
}

func (c *Controller) InsertProduct(ctx context.Context, productId int, product *products.Product) error {
	return c.db.InsertProduct(ctx, productId, product)
}
