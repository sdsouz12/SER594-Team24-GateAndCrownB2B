package catalog

import "context"

type IRepository interface {
	ListProducts(ctx context.Context, category string) ([]*Product, error)
	GetProductById(ctx context.Context, productId int64) (*Product, error)
}
