package orders

import "time"

type Order struct {
	OrderId     int64     `json:"orderId"`
	UserId      int64     `json:"userId"`
	ProductId   int64     `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	Notes       string    `json:"notes,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateOrderRequest struct {
	ProductId int    `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Notes     string `json:"notes"`
}

type CreateOrderResponse struct {
	OrderId   int64  `json:"orderId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
