package database

import "errors"

var ErrNotFound = errors.New("Product not found.")

var ProductAlreadyExists = errors.New("A product with this id already exists.")
