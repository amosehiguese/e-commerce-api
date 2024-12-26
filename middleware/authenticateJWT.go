package middleware

import (
	"net/http"
	"time"

	"github.com/amosehiguese/ecommerce-api/pkg/auth"
	"github.com/gin-gonic/gin"
)

// JWTProtected is the middleware function that validates JWT tokens and processes them.
func JWTProtected() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the access token and refresh token from the cookies
		accessToken, err := c.Cookie("access")
		if err != nil {
			// Handle missing or malformed access token in cookies
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "Missing or malformed access token",
			})
			c.Abort()
			return
		}

		refreshToken, err := c.Cookie("refresh")
		if err != nil {
			// Handle missing or malformed refresh token in cookies
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "Missing or malformed refresh token",
			})
			c.Abort()
			return
		}

		// Validate the access token
		claims, err := auth.ExtractTokenMetadata(c, accessToken)
		if err != nil {
			// If the access token is invalid, return an error
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": true,
				"msg":   "Invalid or expired access token",
			})
			c.Abort()
			return
		}

		// Now we need to check if the access token has expired
		now := time.Now().Unix()
		expAccessToken := claims.Exp
		if now > expAccessToken {
			// Token expired, we need to use the refresh token to get a new access token
			// Parse the refresh token to check its validity
			refreshTokenExp, err := auth.ParseRefreshToken(refreshToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": true,
					"msg":   "Error parsing refresh token",
				})
				c.Abort()
				return
			}

			// If the refresh token is expired, reject the request
			if now > refreshTokenExp {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": true,
					"msg":   "Unauthorized, your session has expired",
				})
				c.Abort()
				return
			}

			// Generate a new access token using the user ID and role
			role, err := auth.VerifyRole(claims.Role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": true,
				})
				return
			}

			newAccessToken, err := auth.GenerateAccessToken(claims.UserID.String(), role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": true,
					"msg":   err.Error(),
				})
				c.Abort()
				return
			}

			// Attach the new access token with its refresh token to the cookie
			auth.AttachToCookie(c, newAccessToken, refreshToken)

			// Proceed with the request after refreshing the token
			c.Next()
			return
		}

		// If the token is still valid, proceed with the request
		c.Set("user", claims.UserID)
		c.Next()
	}
}
