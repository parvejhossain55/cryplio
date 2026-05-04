package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminRoleMiddleware checks if the user has admin role
func AdminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, check if username is "admin" (from seeder)
		// TODO: Add proper role field to user model
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if username.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
