package middleware

import (
	"net/http"

	"github.com/amosehiguese/ecommerce-api/pkg/auth"
	"github.com/gin-gonic/gin"
)

// Admin-only middleware
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token metadata
		tokenMetadata, err := auth.ExtractTokenMetadata(c, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the user is an admin
		if tokenMetadata.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied, admin role required"})
			c.Abort()
			return
		}

		// Proceed if the user is an admin
		c.Next()
	}
}
