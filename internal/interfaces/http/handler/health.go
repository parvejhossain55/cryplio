package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingFunc is a function that checks whether a backing dependency (e.g. the
// database) is reachable. It is injected at construction time so the handler
// itself never imports infrastructure packages.
type PingFunc func(ctx context.Context) error

// HealthHandler handles health and readiness probes.
type HealthHandler struct {
	ping PingFunc // optional; if nil, readiness always reports healthy
}

// NewHealthHandler creates a HealthHandler.
// Pass an optional PingFunc to enable a real readiness check.
func NewHealthHandler(opts ...PingFunc) *HealthHandler {
	h := &HealthHandler{}
	if len(opts) > 0 {
		h.ping = opts[0]
	}
	return h
}

// HealthCheck is a lightweight liveness probe.
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "cryplio-api"})
}

// Liveness reports whether the process is alive.
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alive": true})
}

// Readiness reports whether the service can handle requests.
// If a PingFunc was provided, it is called to verify the database connection.
func (h *HealthHandler) Readiness(c *gin.Context) {
	if h.ping != nil {
		if err := h.ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"ready": false,
				"error": "database unavailable",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"ready": true})
}
