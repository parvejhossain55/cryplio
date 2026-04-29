package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
)

func getLimiter(ip string) *rate.Limiter {
	limiter, exists := limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Second), 100)
		limiters[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
