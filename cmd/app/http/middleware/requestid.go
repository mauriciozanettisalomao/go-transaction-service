package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIdMiddleware is a middleware that adds a request id to the response header
func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Set("requestID", requestID)
		c.Next()
	}
}
