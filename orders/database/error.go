package database

import "errors"

var ErrNotFound = errors.New("Order not found.")

var OrderAlreadyExists = errors.New("An order with this id already exists.")

var UserNotFound = errors.New("The user specified was not found.")

var ProductNotFound = errors.New("The product was not found.")

var ProductAlreadyExists = errors.New("A product with this id already exists.")

var UserAlreadyExists = errors.New("A user with this id already exists.")

var UserOrProductNotFound = errors.New("The user or product specified was not found.")
