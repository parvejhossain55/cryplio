package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	redisstore "cryplio/internal/infrastructure/persistence/redis"

	"github.com/gin-gonic/gin"
)

// NOTE: The in-memory limiter (defaultLimiters) is process-local and will NOT
// work correctly in horizontally-scaled deployments. Pass a RedisRateLimiter to
// RateLimitMiddleware() for distributed, accurate rate limiting.

// ─── In-Memory Fallback ──────────────────────────────────────────────────────

type slidingWindowLimiter struct {
	mu      sync.Mutex
	records map[string][]time.Time
	max     int
	window  time.Duration
}

type rateLimiters struct {
	global   *slidingWindowLimiter
	twoFA    *slidingWindowLimiter
	login    *slidingWindowLimiter
	register *slidingWindowLimiter
	password *slidingWindowLimiter
	email    *slidingWindowLimiter
	session  *slidingWindowLimiter
}

func newSlidingWindowLimiter(max int, window time.Duration) *slidingWindowLimiter {
	return &slidingWindowLimiter{
		records: make(map[string][]time.Time),
		max:     max,
		window:  window,
	}
}

func (rl *slidingWindowLimiter) allow(key string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := now.Add(-rl.window)
	prev := rl.records[key]
	valid := prev[:0]
	for _, t := range prev {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	if len(valid) >= rl.max {
		rl.records[key] = valid
		return false
	}
	rl.records[key] = append(valid, now)
	return true
}

func (rl *slidingWindowLimiter) cleanup(cutoff time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for key, attempts := range rl.records {
		valid := attempts[:0]
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
}

var defaultLimiters = func() *rateLimiters {
	rl := &rateLimiters{
		global:   newSlidingWindowLimiter(100, time.Second),
		twoFA:    newSlidingWindowLimiter(5, 5*time.Minute),
		login:    newSlidingWindowLimiter(10, 5*time.Minute),
		register: newSlidingWindowLimiter(5, 10*time.Minute),
		password: newSlidingWindowLimiter(3, 15*time.Minute),
		email:    newSlidingWindowLimiter(5, 30*time.Minute),
		session:  newSlidingWindowLimiter(50, time.Minute),
	}
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-30 * time.Minute)
			rl.global.cleanup(cutoff)
			rl.twoFA.cleanup(cutoff)
			rl.login.cleanup(cutoff)
			rl.register.cleanup(cutoff)
			rl.password.cleanup(cutoff)
			rl.email.cleanup(cutoff)
			rl.session.cleanup(cutoff)
		}
	}()
	return rl
}()

// ─── Endpoint config ─────────────────────────────────────────────────────────

type endpointRule struct {
	max    int
	window time.Duration
	key    func(c *gin.Context) string
}

var endpointRules = []struct {
	path   string
	method string
	endpointRule
}{
	{"/api/v1/auth/2fa/complete-login", "POST", endpointRule{5, 5 * time.Minute, func(c *gin.Context) string { return "2fa:" + getIP(c) }}},
	{"/api/v1/auth/login", "POST", endpointRule{10, 5 * time.Minute, func(c *gin.Context) string { return "login:" + getIP(c) }}},
	{"/api/v1/auth/register", "POST", endpointRule{5, 10 * time.Minute, func(c *gin.Context) string { return "register:" + getIP(c) }}},
	{"/api/v1/auth/password/reset-request", "POST", endpointRule{3, 15 * time.Minute, func(c *gin.Context) string { return "pwdreset:" + getIP(c) }}},
	{"/api/v1/auth/email/request", "POST", endpointRule{5, 30 * time.Minute, func(c *gin.Context) string { return "email:" + getIP(c) }}},
	{"/api/v1/sessions", "GET", endpointRule{50, time.Minute, func(c *gin.Context) string { return "sessions:" + getUserID(c) }}},
}

// getInMemoryLimiterAndKey looks up the in-memory limiter and rate-limit key
// for the current request path+method.
func getInMemoryLimiterAndKey(c *gin.Context) (*slidingWindowLimiter, string) {
	path := c.Request.URL.Path
	method := c.Request.Method
	for _, rule := range endpointRules {
		if rule.path == path && rule.method == method {
			// map back to the per-endpoint in-memory limiter
			switch rule.path {
			case "/api/v1/auth/2fa/complete-login":
				return defaultLimiters.twoFA, rule.endpointRule.key(c)
			case "/api/v1/auth/login":
				return defaultLimiters.login, rule.endpointRule.key(c)
			case "/api/v1/auth/register":
				return defaultLimiters.register, rule.endpointRule.key(c)
			case "/api/v1/auth/password/reset-request":
				return defaultLimiters.password, rule.endpointRule.key(c)
			case "/api/v1/auth/email/request":
				return defaultLimiters.email, rule.endpointRule.key(c)
			case "/api/v1/sessions":
				return defaultLimiters.session, rule.endpointRule.key(c)
			}
		}
	}
	return defaultLimiters.global, "global:" + getIP(c) + ":" + path
}

func getIP(c *gin.Context) string { return c.ClientIP() }
func getUserID(c *gin.Context) string {
	if uid, exists := c.Get("user_id"); exists {
		return uid.(string)
	}
	return ""
}

// ─── Middleware constructors ──────────────────────────────────────────────────

// RateLimitMiddleware returns an endpoint-aware rate limiting handler.
// If limiter is non-nil (e.g. a RedisRateLimiter), it is used for all checks
// and is safe for horizontally-scaled deployments.
// If limiter is nil, the process-local in-memory sliding-window limiter is used.
func RateLimitMiddleware(limiter ...redisstore.RateLimiter) gin.HandlerFunc {
	var rl redisstore.RateLimiter
	if len(limiter) > 0 && limiter[0] != nil {
		rl = limiter[0]
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		if rl != nil {
			// ── Redis-backed path ────────────────────────────────────────────
			var rule *endpointRule
			for i := range endpointRules {
				if endpointRules[i].path == path && endpointRules[i].method == method {
					rule = &endpointRules[i].endpointRule
					break
				}
			}
			var key string
			var max int
			var window time.Duration
			if rule != nil {
				key = rule.key(c)
				max = rule.max
				window = rule.window
			} else {
				key = "global:" + getIP(c) + ":" + path
				max = 100
				window = time.Second
			}

			allowed, err := rl.Allow(context.Background(), key, max, window)
			if err != nil {
				// Log but fail open on Redis errors
				c.Next()
				return
			}
			if !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many attempts. Please try again later."})
				c.Abort()
				return
			}
		} else {
			// ── In-memory fallback path ──────────────────────────────────────
			inMemLimiter, key := getInMemoryLimiterAndKey(c)
			if !inMemLimiter.allow(key) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many attempts. Please try again later."})
				c.Abort()
				return
			}
		}

		_ = method // suppress unused variable if only path is used above
		c.Next()
	}
}
