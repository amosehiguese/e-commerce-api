package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin-only middleware
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is an admin
		isAdmin := true // Replace with actual logic
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
