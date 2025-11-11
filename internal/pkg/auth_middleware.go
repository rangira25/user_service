package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rangira25/user_service/internal/config"
)

// JWTMiddleware verifies the JWT and extracts user data from it.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		// 2. Check for "Bearer " prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization format"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Parse and validate the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Ensure token method is HMAC (to prevent tampering)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// 4. Extract claims safely
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// 5. Extract user data
		var userID, role string

		// Handle user_id as string or number
		if v, ok := claims["user_id"]; ok {
			switch val := v.(type) {
			case string:
				userID = val
			case float64:
				userID = fmt.Sprintf("%.0f", val)
			}
		}

		if v, ok := claims["role"].(string); ok {
			role = v
		}

		if userID == "" || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user data in token"})
			return
		}

		// 6. Store values in Gin context for later use
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}

// RequireRole ensures only specific roles can access certain routes.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no role found in context"})
			return
		}

		role, ok := roleVal.(string)
		if !ok || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
			return
		}

		c.Next()
	}
}
