package aiproxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	aiServiceURL string
}

func NewHandler(aiServiceURL string) *Handler {
	return &Handler{aiServiceURL: aiServiceURL}
}

// Proxy forwards the request body to the AI service and returns its response.
func (h *Handler) Proxy(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request"})
			return
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, h.aiServiceURL+path, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build upstream request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI service unavailable"})
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, "application/json", respBody)
	}
}
