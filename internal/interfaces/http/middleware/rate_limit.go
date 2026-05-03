package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// endpointRateLimiter implements sliding window rate limiting per key
type endpointRateLimiter struct {
	mu      sync.RWMutex
	records map[string][]time.Time
	max     int
	window  time.Duration
}

func newEndpointRateLimiter(max int, window time.Duration) *endpointRateLimiter {
	rl := &endpointRateLimiter{
		records: make(map[string][]time.Time),
		max:     max,
		window:  window,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *endpointRateLimiter) allow(key string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	attempts := rl.records[key]
	cutoff := now.Add(-rl.window)
	valid := make([]time.Time, 0, rl.max)
	for _, t := range attempts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	if len(valid) >= rl.max {
		rl.records[key] = valid
		return false
	}
	valid = append(valid, now)
	rl.records[key] = valid
	return true
}

func (rl *endpointRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-rl.window)
		for key, attempts := range rl.records {
			valid := make([]time.Time, 0, len(attempts))
			for _, t := range attempts {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.records, key)
			} else {
				rl.records[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

// Global limiters per endpoint
var (
	globalLimiter    *endpointRateLimiter
	twoFALimiter     *endpointRateLimiter
	loginLimiter     *endpointRateLimiter
	registerLimiter  *endpointRateLimiter
	passwordLimiter  *endpointRateLimiter
	emailLimiter     *endpointRateLimiter
	sessionLimiter   *endpointRateLimiter
)

func init() {
	globalLimiter   = newEndpointRateLimiter(100, time.Second)
	twoFALimiter    = newEndpointRateLimiter(5, 5*time.Minute)
	loginLimiter    = newEndpointRateLimiter(10, 5*time.Minute)
	registerLimiter = newEndpointRateLimiter(5, 10*time.Minute)
	passwordLimiter = newEndpointRateLimiter(3, 15*time.Minute)
	emailLimiter    = newEndpointRateLimiter(5, 30*time.Minute)
	sessionLimiter  = newEndpointRateLimiter(50, time.Minute)
}

func getIP(c *gin.Context) string {
	return c.ClientIP()
}

func getUserID(c *gin.Context) string {
	uid, exists := c.Get("user_id")
	if exists {
		return uid.(string)
	}
	return ""
}

// RateLimitMiddleware applies rate limiting based on endpoint
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		var limiter *endpointRateLimiter
		var key string

		switch {
		case path == "/api/v1/auth/2fa/complete-login" && method == "POST":
			limiter = twoFALimiter
			key = "2fa:" + getIP(c)
		case path == "/api/v1/auth/login" && method == "POST":
			limiter = loginLimiter
			key = "login:" + getIP(c)
		case path == "/api/v1/auth/register" && method == "POST":
			limiter = registerLimiter
			key = "register:" + getIP(c)
		case path == "/api/v1/auth/password/reset-request" && method == "POST":
			limiter = passwordLimiter
			key = "pwdreset:" + getIP(c)
		case path == "/api/v1/auth/email/request" && method == "POST":
			limiter = emailLimiter
			key = "email:" + getIP(c)
		case path == "/api/v1/sessions" && method == "GET":
			limiter = sessionLimiter
			key = "sessions:" + getUserID(c)
		default:
			limiter = globalLimiter
			key = "global:" + getIP(c) + ":" + path
		}

		if key == "" {
			key = getIP(c)
		}

		if !limiter.allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
