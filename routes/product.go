package routes

import (
	"github.com/amosehiguese/ecommerce-api/api"
	"github.com/amosehiguese/ecommerce-api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(router *gin.RouterGroup, a api.API) {
	// Admin routes for product management
	admin := router.Group("/products", middleware.AdminOnly())
	{
		admin.POST("/", a.CreateProduct)
		admin.PUT("/:id", a.UpdateProduct)
		admin.DELETE("/:id", a.DeleteProduct)
	}

	// General product routes
	router.GET("/products", a.ListProducts)
	router.GET("/products/:id", a.GetProduct)
}
