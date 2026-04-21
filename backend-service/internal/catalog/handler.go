package catalog

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{service: service}
}

// ListProducts handles GET /api/catalog/products
func (h *Handler) ListProducts(c *gin.Context) {
	category := c.Query("category")
	products, err := h.service.ListProducts(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	if products == nil {
		products = []*Product{}
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

// SemanticSearch handles POST /api/catalog/search
func (h *Handler) SemanticSearch(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	results, err := h.service.SemanticSearch(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Search unavailable: " + err.Error()})
		return
	}
	if results == nil {
		results = []*SearchResult{}
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}
