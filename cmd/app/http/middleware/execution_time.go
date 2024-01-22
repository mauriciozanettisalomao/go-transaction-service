package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// ExecutionTimeMiddleware is a middleware that logs the execution time of the request
func ExecutionTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsReq := time.Now()

		c.Next()

		execTime := time.Since(tsReq).Seconds()

		slog.Info("execution time",
			"path", c.FullPath(),
			"execution_time_seconds", execTime,
		)
	}
}
