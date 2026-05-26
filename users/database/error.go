package database

import "errors"

var ErrNotFound = errors.New("User not found.")

var UserAlreadyExists = errors.New("A user with this id already exists.")
