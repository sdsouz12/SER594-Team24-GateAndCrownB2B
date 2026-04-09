package auth

import (
	"errors"
	"net/http"

	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/internal/ctxutil"
	"github.com/gin-gonic/gin"
)

// Handler handles auth HTTP requests
type Handler struct {
	service IService
}

// NewHandler creates a new auth handler
func NewHandler(service IService) *Handler {
	return &Handler{service: service}
}

// Login handles POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	response, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
		if err == ErrAccountDeactivated {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account has been deactivated",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMe handles GET /auth/me
func (h *Handler) GetMe(c *gin.Context) {
	caller, ok := ctxutil.GetCaller(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	me, err := h.service.GetMe(c.Request.Context(), caller.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": me})
}

// UpdateMe handles PATCH /auth/me — profile and/or password for the authenticated user.
func (h *Handler) UpdateMe(c *gin.Context) {
	caller, ok := ctxutil.GetCaller(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	me, err := h.service.UpdateMe(c.Request.Context(), caller.UserId, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrNothingToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrPasswordTooShort), errors.Is(err, ErrPasswordChangeRequiresCurrent):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrWrongCurrentPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": me})
}
