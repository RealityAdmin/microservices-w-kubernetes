package memory

import (
	"context"
	"sync"

	"k-microserv-kuber.com/orders"
	"k-microserv-kuber.com/orders/database"
)

type MemoryDB struct {
	sync.RWMutex
	data map[int]*orders.Order
}

func New() *MemoryDB {
	return &MemoryDB{data: make(map[int]*orders.Order)}
}

// Get the order by id
func (d *MemoryDB) GetOrder(_ context.Context, orderId int) (*orders.Order, error) {
	d.RLock()
	defer d.RUnlock()

	order, ok := d.data[orderId]
	if !ok {
		return nil, database.ErrNotFound
	}

	return order, nil
}

// Place an order. Assume that user and product exists (checks done already)
func (d *MemoryDB) PlaceOrder(_ context.Context, userId int, productId int) error {
	d.RLock()
	defer d.RUnlock()

	new_id := len(d.data)
	d.data[new_id] = &orders.Order{
		OrderID:   new_id,
		UserID:    userId,
		ProductID: productId,
	}

	return nil
}
