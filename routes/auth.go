package routes

import (
	"github.com/amosehiguese/ecommerce-api/api"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, a api.API) {
	router.POST("/register", a.Register)
	router.POST("/login", a.Login)

	router.POST("/create-admin", a.CreateAdmin)
}

func RegisterTokenRenewal(router *gin.RouterGroup, a api.API) {
	router.POST("/renew-token", a.RenewTokens)
}
