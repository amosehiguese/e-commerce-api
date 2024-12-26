package routes

import (
	"github.com/amosehiguese/ecommerce-api/api"
	"github.com/amosehiguese/ecommerce-api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.RouterGroup, a api.API) {
	// Order routes for all authenticated users
	router.POST("/orders", a.CreateOrder)
	router.GET("/orders", a.ListUserOrders)
	router.PUT("/orders/:id/cancel", a.CancelOrder)

	// Admin-only order status management
	admin := router.Group("/orders", middleware.AdminOnly())
	{
		admin.PUT("/:id/status", a.UpdateOrderStatus)
	}
}
