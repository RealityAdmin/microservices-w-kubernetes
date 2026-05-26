package memory

import (
	"context"
	"sync"

	"k-microserv-kuber.com/users"
	"k-microserv-kuber.com/users/database"
)

type MemoryDB struct {
	sync.RWMutex
	data map[int]*users.User
}

func New() *MemoryDB {
	return &MemoryDB{data: make(map[int]*users.User)}
}

func NewTest() *MemoryDB {
	return &MemoryDB{data: map[int]*users.User{1: {}}}
}

func (d *MemoryDB) GetUser(_ context.Context, id int) (*users.User, error) {
	d.RLock()
	defer d.RUnlock()

	// Find the user
	user, ok := d.data[id]
	if !ok {
		return nil, database.ErrNotFound
	}
	return user, nil
}

func (d *MemoryDB) PutUser(_ context.Context, id int, user *users.User) error {
	d.Lock()
	defer d.Unlock()

	// Assert that the user does not already exist.
	_, ok := d.data[id]
	if ok {
		return database.UserAlreadyExists
	}

	// Insert the user
	d.data[id] = user

	return nil
}
