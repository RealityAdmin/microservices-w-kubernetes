package memory

import (
	"context"
	"sync"

	"k-microserv-kuber.com/products"
	"k-microserv-kuber.com/products/database"
)

type MemoryDB struct {
	sync.RWMutex
	data map[int]*products.Product
}

func New() *MemoryDB {
	return &MemoryDB{data: make(map[int]*products.Product)}
}

// Find a product by Id
func (d *MemoryDB) GetProduct(_ context.Context, productId int) (*products.Product, error) {
	d.RLock()
	defer d.RUnlock()

	product, ok := d.data[productId]
	if !ok {
		return nil, database.ErrNotFound
	}

	return product, nil
}

// Insert a product if the id is not taken
func (d *MemoryDB) InsertProduct(_ context.Context, productId int, product *products.Product) error {

	d.Lock()
	defer d.Unlock()

	_, ok := d.data[productId]
	if ok {
		return database.ProductAlreadyExists
	}

	d.data[productId] = product

	return nil
}
