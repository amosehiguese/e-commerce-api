package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register a new user
func (api *API) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login user and issue a JWT
func (api *API) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "JWT_TOKEN_HERE"})
}
