package orders

import "context"

type Service struct {
	repo IRepository
}

func NewService(repo IRepository) IService {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, userId int64, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	orderId, err := s.repo.CreateOrder(ctx, userId, req.ProductId, req.Quantity, req.Notes)
	if err != nil {
		return nil, err
	}
	return &CreateOrderResponse{
		OrderId: orderId,
		Status:  "pending",
		Message: "Your order has been received! Our team will review it and get back to you shortly.",
	}, nil
}

func (s *Service) GetMyOrders(ctx context.Context, userId int64) ([]Order, error) {
	return s.repo.GetOrdersByUserId(ctx, userId)
}
