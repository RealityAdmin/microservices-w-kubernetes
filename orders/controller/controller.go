package controller

import (
	"context"
	"errors"

	"k-microserv-kuber.com/orders"
	"k-microserv-kuber.com/orders/database"
	"k-microserv-kuber.com/orders/gateway"
	"k-microserv-kuber.com/products"
	"k-microserv-kuber.com/users"
)

type Database interface {
	GetOrder(ctx context.Context, orderId int) (*orders.Order, error)
	PlaceOrder(ctx context.Context, orderId int, order *orders.Order) error
}

type Controller struct {
	db             Database
	userGateway    userGateway
	productGateway productGateway
}

type userGateway interface {
	GetUser(ctx context.Context, userId int) (*users.User, error)
	PutUser(ctx context.Context, userId int, email string) error
}

type productGateway interface {
	GetProduct(ctx context.Context, productId string) (*products.Product, error)
	InsertProduct(ctx context.Context, productId string, productName string, price float64) error
}

func New(db Database, ug userGateway, pg productGateway) *Controller {
	return &Controller{db: db, userGateway: ug, productGateway: pg}
}

func (c *Controller) GetUser(ctx context.Context, userId int) (*users.User, error) {
	u, err := c.userGateway.GetUser(ctx, userId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, gateway.ErrNotFound
	}

	return u, err
}

// TODO: Integrate the user and product gateway calls into this.
func (c *Controller) GetOrder(ctx context.Context, orderId int) (*orders.Order, error) {
	o, err := c.db.GetOrder(ctx, orderId)
	if err != nil && errors.Is(err, database.ErrNotFound) {
		return nil, database.ErrNotFound
	}
	return o, err
}

func (c *Controller) PlaceOrder(ctx context.Context, orderId int, order *orders.Order) error {
	return c.db.PlaceOrder(ctx, orderId, order)
}
