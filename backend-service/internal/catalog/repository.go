package catalog

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

func (r *Repository) ListProducts(ctx context.Context, category string) ([]*Product, error) {
	query := `
		SELECT product_id, name, category, description, price_range, material, dimensions, status, created_at, updated_at
		FROM catalog_product
		WHERE status = 'active'
	`
	args := []any{}
	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}
	query += " ORDER BY product_id"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ProductId, &p.Name, &p.Category, &p.Description,
			&p.PriceRange, &p.Material, &p.Dimensions, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, &p)
	}
	return products, rows.Err()
//Pin functionality here is to get the product details by its ID. This will be used when user clicks on a product to view its details.

func (r *Repository) GetProductById(ctx context.Context, productId int64) (*Product, error) {
	var p Product
	err := r.db.QueryRow(ctx, `
		SELECT product_id, name, category, description, price_range, material, dimensions, status, created_at, updated_at
		FROM catalog_product WHERE product_id = $1
	`, productId).Scan(&p.ProductId, &p.Name, &p.Category, &p.Description,
		&p.PriceRange, &p.Material, &p.Dimensions, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return &p, nil
}

func (r *Repository) countProducts(ctx context.Context, category string) (int, error) {
	query := SELECT COUNT(*) FROM catalog_product WHERE status = 'active'
	args := []any{}
	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}
	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}
	return count, nil
}
