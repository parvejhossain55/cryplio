package middleware

import (
	"net/http"

	sharedjwt "cryplio/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT from cookie or Authorization header
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := sharedjwt.FromRequest(readAuthCookie(c), c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, err := sharedjwt.Parse(secret, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		userID, ok := claims[sharedjwt.ClaimUserID].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			c.Abort()
			return
		}

		username, _ := claims["username"].(string)

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("token_type", claims[sharedjwt.ClaimTokenType])
		if jti, ok := claims["jti"].(string); ok {
			c.Set("token_id", jti)
		}
		c.Next()
	}
}

// OptionalAuth allows endpoints that work with or without auth
func OptionalAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := sharedjwt.FromRequest(readAuthCookie(c), c.GetHeader("Authorization"))
		if err != nil {
			c.Next()
			return
		}

		if tokenString != "" {
			claims, err := sharedjwt.Parse(secret, tokenString)
			if err == nil {
				if userID, ok := claims[sharedjwt.ClaimUserID].(string); ok {
					username, _ := claims["username"].(string)
					c.Set("user_id", userID)
					c.Set("username", username)
					c.Set("token_type", claims[sharedjwt.ClaimTokenType])
					if jti, ok := claims["jti"].(string); ok {
						c.Set("token_id", jti)
					}
				}
			}
		}
		c.Next()
	}
}

func readAuthCookie(c *gin.Context) string {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return ""
	}
	return cookie
}
