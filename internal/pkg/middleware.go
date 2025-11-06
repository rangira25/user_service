package pkg

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs requests with method, path, status, and duration
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[%d] %s %s (%s)", status, c.Request.Method, c.Request.URL.Path, duration)
	}
}

// RecoveryMiddleware can be used to recover from panics and return JSON error
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
