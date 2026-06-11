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
	GetProduct(ctx context.Context, productId int) (*products.Product, error)
	InsertProduct(ctx context.Context, productId int, productName string, price float64) error
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

func (c *Controller) PutUser(ctx context.Context, userId int, email string) error {
	return c.userGateway.PutUser(ctx, userId, email)
}

func (c *Controller) GetProduct(ctx context.Context, productId int) (*products.Product, error) {
	p, err := c.productGateway.GetProduct(ctx, productId)
	return p, err
}

func (c *Controller) InsertProduct(ctx context.Context, productId int, productName string, price float64) error {
	return c.productGateway.InsertProduct(ctx, productId, productName, price)
}

func (c *Controller) GetOrder(ctx context.Context, orderId int) (*orders.Order, error) {
	o, err := c.db.GetOrder(ctx, orderId)
	if err != nil && errors.Is(err, database.ErrNotFound) {
		return nil, database.ErrNotFound
	}
	return o, err
}

// TODO: Integrate the user and product gateway calls into this.
func (c *Controller) PlaceOrder(ctx context.Context, orderId int, order *orders.Order) error {
	productId := order.ProductID
	userId := order.UserID

	// Assert that users and products exist
	// var wg sync.WaitGroup

	// userExists :=
	// var productExists bool

	// wg.Go(func ()  {
	// 	_, err := c.userGateway.GetUser(ctx, userId)
	// 	if err != nil {
	// 		userExists = false
	// 	}
	// 	userExists = true
	// })

	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan bool)

	// Launch subroutine to get user
	go func(context.Context, chan bool) {
		_, err := c.userGateway.GetUser(newCtx, userId)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			errChan <- false
			return
		}
		errChan <- true
	}(newCtx, errChan)

	// Launch subroutine to get product
	go func(context.Context, chan bool) {
		_, err := c.productGateway.GetProduct(newCtx, productId)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			errChan <- false
			return
		}
		errChan <- true
	}(newCtx, errChan)

	// If either fails, cancel the other and return an error.
	for i := 0; i < 2; i++ {
		err := <-errChan
		if err == false {
			cancel()
			return database.UserOrProductNotFound
		}
	}

	// We have a match, now actually place the order
	return c.db.PlaceOrder(ctx, orderId, order)
}
