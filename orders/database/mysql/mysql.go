package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"k-microserv-kuber.com/orders"
	"k-microserv-kuber.com/orders/database"
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

func (r *Repository) GetOrder(ctx context.Context, orderId int) (*orders.Order, error) {
	var id int
	var userId int
	var productId int

	row := r.db.QueryRowContext(ctx, "SELECT id, userId, productId FROM Orders WHERE id = ?", orderId)
	if err := row.Scan(&id, &userId, &productId); err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &orders.Order{
		OrderID:   id,
		UserID:    userId,
		ProductID: productId,
	}, nil
}

func (r *Repository) PlaceOrder(ctx context.Context, orderId int, order *orders.Order) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO Orders (userId, productId) VALUES(?, ?)", order.UserID, order.ProductID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected < 1 {
		return database.OrderAlreadyExists
	}

	return nil
}
