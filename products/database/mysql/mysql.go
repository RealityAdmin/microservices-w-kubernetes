package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"k-microserv-kuber.com/products"
	"k-microserv-kuber.com/products/database"
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

// Get a product by id
func (r *Repository) GetProduct(ctx context.Context, productId int) (*products.Product, error) {
	var id int
	var name string
	var price float64

	row := r.db.QueryRowContext(ctx, "SELECT id, name, price FROM Products WHERE id = ?", productId)
	if err := row.Scan(&id, &name, &price); err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &products.Product{
		ProductID:   id,
		ProductName: name,
		Price:       price,
	}, nil
}

// Insert a product into the database
func (r *Repository) InsertProduct(ctx context.Context, productId int, product *products.Product) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO Products (name, price) VALUES(?, ?)", product.ProductName, product.Price)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return database.ProductAlreadyExists
	}

	return nil
}
