package orders

import "context"

type IRepository interface {
	CreateOrder(ctx context.Context, userId int64, productId, quantity int, notes string) (int64, error)
	GetOrdersByUserId(ctx context.Context, userId int64) ([]Order, error)
}
