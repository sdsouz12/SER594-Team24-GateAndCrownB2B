package orders

import "context"

type IService interface {
	CreateOrder(ctx context.Context, userId int64, req *CreateOrderRequest) (*CreateOrderResponse, error)
	GetMyOrders(ctx context.Context, userId int64) ([]Order, error)
}
