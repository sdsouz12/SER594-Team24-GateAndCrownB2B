package orders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) IRepository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrder(ctx context.Context, userId int64, productId, quantity int, notes string) (int64, error) {
	var orderId int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO orders (user_id, product_id, quantity, notes, status)
		VALUES ($1, $2, $3, $4, 'pending')
		RETURNING order_id
	`, userId, productId, quantity, notes).Scan(&orderId)
	if err != nil {
		return 0, fmt.Errorf("create order: %w", err)
	}
	return orderId, nil
}

func (r *Repository) GetOrdersByUserId(ctx context.Context, userId int64) ([]Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT o.order_id, o.user_id, o.product_id, p.name, o.quantity, COALESCE(o.notes, ''), o.status, o.created_at
		FROM orders o
		JOIN catalog_product p ON p.product_id = o.product_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}
	defer rows.Close()

	var result []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderId, &o.UserId, &o.ProductId, &o.ProductName, &o.Quantity, &o.Notes, &o.Status, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		result = append(result, o)
	}
	return result, rows.Err()
}
