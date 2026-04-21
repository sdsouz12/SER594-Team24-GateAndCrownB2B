package catalog

import "time"

type Product struct {
	ProductId   int64     `json:"productId"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description,omitempty"`
	PriceRange  string    `json:"priceRange,omitempty"`
	Material    string    `json:"material,omitempty"`
	Dimensions  string    `json:"dimensions,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SearchRequest struct {
	Query    string `json:"query" binding:"required"`
	Category string `json:"category"`
	Limit    int    `json:"limit"`
}

type SearchResult struct {
	Product    Product `json:"product"`
	Similarity float64 `json:"similarity"`
}
