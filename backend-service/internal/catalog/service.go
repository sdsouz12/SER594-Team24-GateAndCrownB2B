package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service struct {
	repo         IRepository
	aiServiceURL string
}

func NewService(repo IRepository, aiServiceURL string) IService {
	return &Service{repo: repo, aiServiceURL: aiServiceURL}
}

func (s *Service) ListProducts(ctx context.Context, category string) ([]*Product, error) {
	return s.repo.ListProducts(ctx, category)
}

//Nurs code here

type aiSearchRequest struct {
	Query    string `json:"query"`
	Category string `json:"category,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type aiSearchItem struct {
	ProductId  int64   `json:"product_id"`
	Similarity float64 `json:"similarity"`
}

type aiSearchResponse struct {
	Results []aiSearchItem `json:"results"`
}

func (s *Service) SemanticSearch(ctx context.Context, req *SearchRequest) ([]*SearchResult, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	body, _ := json.Marshal(aiSearchRequest{Query: req.Query, Category: req.Category, Limit: limit})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.aiServiceURL+"/search", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build ai request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ai service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai service error: status %d", resp.StatusCode)
	}

	var aiResp aiSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("decode ai response: %w", err)
	}

	var results []*SearchResult
	for _, item := range aiResp.Results {
		product, err := s.repo.GetProductById(ctx, item.ProductId)
		if err != nil {
			continue
		}
		results = append(results, &SearchResult{Product: *product, Similarity: item.Similarity})
	}
	return results, nil
}
