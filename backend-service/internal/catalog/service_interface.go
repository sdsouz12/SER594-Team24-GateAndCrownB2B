package catalog

import "context"

type IService interface {
	ListProducts(ctx context.Context, category string) ([]*Product, error)
	SemanticSearch(ctx context.Context, req *SearchRequest) ([]*SearchResult, error)
}
