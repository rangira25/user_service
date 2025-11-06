package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole ensures a user has a specific role (e.g., "admin")
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-Role") // simulate auth via request header for demo
		if role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: admin only"})
			return
		}
		c.Next()
	}
}
