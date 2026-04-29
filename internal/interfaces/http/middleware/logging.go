package middleware

import (
	"time"

	"cryplio/pkg/logger"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		logFields := logger.Fields{
			"method":     method,
			"status":     status,
			"latency":    latency.String(),
			"path":       path,
			"client_ip":  clientIP,
			"user_agent": c.Request.UserAgent(),
		}
		if query != "" {
			logFields["query"] = query
		}
		if userID, exists := c.Get("user_id"); exists {
			logFields["user_id"] = userID
		}

		if status >= 500 {
			logger.Error("request failed", logFields)
		} else if status >= 400 {
			logger.Warn("request error", logFields)
		} else {
			logger.Info("request completed", logFields)
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(500, gin.H{"error": err})
		}
		c.AbortWithStatus(500)
	})
}
