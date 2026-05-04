package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/ctxutil"
)

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	caller, ok := ctxutil.GetCaller(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	resp, err := h.service.CreateOrder(c.Request.Context(), caller.UserId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *Handler) GetMyOrders(c *gin.Context) {
	caller, ok := ctxutil.GetCaller(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orders, err := h.service.GetMyOrders(c.Request.Context(), caller.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	if orders == nil {
		orders = []Order{}
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}
