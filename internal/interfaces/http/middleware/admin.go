package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminRoleMiddleware checks if the user has admin role
func AdminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
