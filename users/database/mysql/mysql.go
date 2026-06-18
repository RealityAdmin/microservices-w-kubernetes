package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"k-microserv-kuber.com/users"
	"k-microserv-kuber.com/users/database"
)

type Repository struct {
	db *sql.DB
}

// Create new connection to MySQL
func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@tcp(mysql-external:3306)/orders")
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, err
}

// Get a user from the database by id
func (r *Repository) GetUser(ctx context.Context, userId int) (*users.User, error) {
	var id int
	var email string

	row := r.db.QueryRowContext(ctx, "SELECT id, email FROM Users WHERE id = ?", userId)
	if err := row.Scan(&id, &email); err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &users.User{
		UserId: id,
		Email:  email,
	}, nil
}

// Insert a user into the database
func (r *Repository) PutUser(ctx context.Context, userId int, user *users.User) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO Users (email) VALUES(?)", user.Email)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return database.UserAlreadyExists
	}

	return nil
}
