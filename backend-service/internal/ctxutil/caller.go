package ctxutil

import (
	"github.com/gin-gonic/gin"
)

// CallerInfo holds the authenticated user's identity extracted from the JWT context.
// Available in every handler after AuthMiddleware runs.
type CallerInfo struct {
	UserId         int64
	Username       string
	Role           string
	OrganizationId *int64
}

// GetCaller extracts the logged-in user's info from the Gin context.
// Returns (nil, false) if the context has not been populated by AuthMiddleware.
//
// Usage (any handler, any package):
//
//	caller, ok := middleware.GetCaller(c)
//	if !ok {
//	    c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
//	    return
//	}
//	fmt.Println(caller.UserId, caller.Role)
func GetCaller(c *gin.Context) (*CallerInfo, bool) {
	userIdRaw, ok1 := c.Get("user_id")
	usernameRaw, ok2 := c.Get("username")
	roleRaw, ok3 := c.Get("role")
	if !ok1 || !ok2 || !ok3 {
		return nil, false
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		return nil, false
	}
	username, _ := usernameRaw.(string)
	role, _ := roleRaw.(string)

	var orgId *int64
	if raw, exists := c.Get("organization_id"); exists {
		if oid, ok := raw.(*int64); ok {
			orgId = oid
		}
	}

	return &CallerInfo{
		UserId:         userId,
		Username:       username,
		Role:           role,
		OrganizationId: orgId,
	}, true
}
