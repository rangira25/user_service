package pkg

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs every request with status, method, path, and duration.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[%d] %s %s (%s)", status, c.Request.Method, c.Request.URL.Path, duration)
	}
}

// RecoveryMiddleware recovers from panics and returns a standardized JSON error.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				c.JSON(500, gin.H{
					"success": false,
					"error":   "Internal Server Error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
