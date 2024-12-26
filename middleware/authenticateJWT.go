package middleware

import (
	"github.com/gin-gonic/gin"
)

// JWT authentication middleware
func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement JWT validation logic
		c.Next()
	}
}
