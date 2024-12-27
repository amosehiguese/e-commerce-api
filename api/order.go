package api

import (
	"net/http"
	"time"

	"github.com/amosehiguese/ecommerce-api/api/payload"
	"github.com/amosehiguese/ecommerce-api/pkg/auth"
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/pkg/validator"
	"github.com/amosehiguese/ecommerce-api/query"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CreateOrder godoc
// @Summary      Create an Order
// @Description  Create a new order for the authenticated user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        orderPayload body payload.OrderPayload true "Order Payload"
// @Success      200 {object} map[string]interface{} "Order created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized, token expired"
// @Failure      403 {object} map[string]interface{} "Permission denied"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /api/orders [post]
func (api *API) CreateOrder(c *gin.Context) {
	log := logger.Get()

	claims, err := auth.ExtractTokenMetadata(c)
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	if !claims.Credentials[auth.OrderCreateCredential] {
		log.Warn("Permission denied for order create", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	var orderPayload payload.OrderPayload
	if err := c.ShouldBindJSON(&orderPayload); err != nil {
		log.Error("Invalid JSON for order", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	validate := validator.NewValidator()
	if err := validate.Struct(orderPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	order := &query.Order{
		ID:          uuid.New(),
		UserID:      claims.UserID,
		TotalAmount: orderPayload.CalculateOrderTotal(),
	}

	for _, itemPayload := range orderPayload.Items {
		product, err := api.Q.GetProductByID(c, itemPayload.ProductID)
		itemPrice := decimal.NewFromFloat(itemPayload.Price)

		if itemPayload.Quantity > product.UnitsInStock {
			log.Error("Product quantity greater than units in stock", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "product quantity greater than units in stock: " + err.Error(),
			})
			return
		}

		if itemPrice != product.Price {
			log.Error("price inconsistency", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   "item price differ from product price: " + err.Error(),
			})
			return

		}

		order.Items = append(order.Items, query.OrderItem{
			ProductID: itemPayload.ProductID,
			Quantity:  itemPayload.Quantity,
			Price:     itemPrice,
			CreatedAt: time.Now(),
		})
	}

	if _, err := api.Q.CreateOrder(c, order); err != nil {
		log.Error("Error placing order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Order created successfully", zap.String("order_id", order.ID.String()))
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "Order placed successfully", "order_id": order.ID.String()})
}

// ListUserOrders godoc
// @Summary      List User Orders
// @Description  Retrieve all orders for the authenticated user
// @Tags         Orders
// @Produce      json
// @Success      200 {object} map[string]interface{} "User's orders retrieved successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized, token expired"
// @Failure      403 {object} map[string]interface{} "Permission denied"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /api/orders [get]
func (api *API) ListUserOrders(c *gin.Context) {
	log := logger.Get()

	claims, err := auth.ExtractTokenMetadata(c)
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	if !claims.Credentials[auth.OrderReadCredential] {
		log.Warn("Permission denied for order read", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	orders, err := api.Q.GetOrdersByUserID(c, claims.UserID)
	if err != nil {
		log.Error("Error fetching orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Fetched user orders successfully", zap.Int("count", len(orders)))
	c.JSON(http.StatusOK, gin.H{"error": false, "orders": orders})
}

// CancelOrder godoc
// @Summary      Cancel an Order
// @Description  Cancel a specific order if it is still pending
// @Tags         Orders
// @Param        id path string true "Order ID"
// @Produce      json
// @Success      200 {object} map[string]interface{} "Order cancelled successfully"
// @Failure      401 {object} map[string]interface{} "Unauthorized, token expired"
// @Failure      403 {object} map[string]interface{} "Permission denied"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /api/orders/{id}/cancel [patch]
func (api *API) CancelOrder(c *gin.Context) {
	log := logger.Get()

	claims, err := auth.ExtractTokenMetadata(c)
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	if !claims.Credentials[auth.OrderCancelCredential] {
		log.Warn("Permission denied for update read", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	orderID := uuid.MustParse(c.Param("id"))

	if err := api.Q.CancelOrderIfPending(c, orderID); err != nil {
		log.Error("Error canceling order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Order cancelled successfully", zap.String("order_id", orderID.String()))
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "Order cancelled successfully"})
}

// UpdateOrderStatus godoc
// @Summary      Update Order Status
// @Description  Update the status of a specific order
// @Tags         Orders
// @Param        id path string true "Order ID"
// @Param        orderUpdatePayload body payload.OrderUpdatePayload true "Order Update Payload"
// @Produce      json
// @Success      200 {object} map[string]interface{} "Order status updated successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Permission denied"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /orders/{id}/status [put]
func (api *API) UpdateOrderStatus(c *gin.Context) {
	log := logger.Get()

	claims, err := auth.ExtractTokenMetadata(c)
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if !claims.Credentials[auth.OrderUpdateCredential] {
		log.Warn("Permission denied for update read", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	var orderUpdatePayload payload.OrderUpdatePayload
	if err := c.ShouldBindJSON(&orderUpdatePayload); err != nil {
		log.Error("Invalid JSON for update order status", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	validate := validator.NewValidator()
	if err := validate.Struct(orderUpdatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	orderID := uuid.MustParse(c.Param("id"))
	if err := api.Q.UpdateOrderStatus(c, orderID, orderUpdatePayload.Status); err != nil {
		log.Error("Error updating order status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Order status updated successfully", zap.String("order_id", orderID.String()), zap.String("status", orderUpdatePayload.Status))
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "Order status updated successfully"})
}
