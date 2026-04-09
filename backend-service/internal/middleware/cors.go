package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS handles browser cross-origin requests. allowedOrigins is a comma-separated
// list (e.g. "http://host:80,http://host:8081"); spaces are trimmed.
// If the request Origin header matches one entry exactly, that origin is echoed
// back (required when using credentials). If there is no Origin, the first
// allowed origin is used as a default for non-browser clients.
func CORS(allowedOriginsCSV string) gin.HandlerFunc {
	allowed := parseAllowedOrigins(allowedOriginsCSV)
	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		allow := ""
		if origin != "" {
			for _, o := range allowed {
				if origin == o {
					allow = origin
					break
				}
			}
		} else if len(allowed) > 0 {
			// No Origin: same-origin, curl, etc.
			allow = allowed[0]
		}
		if allow != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allow)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func parseAllowedOrigins(csv string) []string {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
