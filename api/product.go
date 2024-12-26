package api

import (
	"net/http"
	"time"

	"github.com/amosehiguese/ecommerce-api/pkg/auth"
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/query"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product by providing product details
// @Tags products
// @Accept json
// @Produce json
// @Param product body query.Product true "Product data"
// @Security CookieAuth
// @Success 200 {object} query.Product
// @Failure 400 {object} gin.H{"error": true, "msg": "error message"}
// @Failure 401 {object} gin.H{"error": true, "msg": "unauthorized"}
// @Failure 403 {object} gin.H{"error": true, "msg": "permission denied"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router /api/products [post]
func (api *API) CreateProduct(c *gin.Context) {
	log := logger.Get()

	// Extract claims from the token
	claims, err := auth.ExtractTokenMetadata(c, "access")
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Check token expiration
	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	// Check if the user has permission to create a product
	if !claims.Credentials["product:create"] {
		log.Warn("Permission denied for product creation", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	// Bind the JSON request to the product struct
	var product query.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Error("Invalid JSON for product", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Set the ID and created/updated timestamps
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Perform the DB operation to create the product
	if err := api.Q.CreateProduct(c, &product); err != nil {
		log.Error("Error creating product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Product created successfully", zap.String("product_id", product.ID.String()))
	c.JSON(http.StatusOK, gin.H{"error": false, "product": product})
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update a product's details using its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body query.Product true "Updated product data"
// @Security CookieAuth
// @Success 200 {object} query.Product
// @Failure 400 {object} gin.H{"error": true, "msg": "error message"}
// @Failure 401 {object} gin.H{"error": true, "msg": "unauthorized"}
// @Failure 403 {object} gin.H{"error": true, "msg": "permission denied"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router /api/products/{id} [put]
func (api *API) UpdateProduct(c *gin.Context) {
	log := logger.Get()

	// Extract claims from the token
	claims, err := auth.ExtractTokenMetadata(c, "access")
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Check token expiration
	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	// Check if the user has permission to update the product
	if !claims.Credentials["product:update"] {
		log.Warn("Permission denied for product update", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	// Get the product ID from the URL parameter
	productID := c.Param("id")

	// Bind the JSON request to the product struct
	var product query.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Error("Invalid JSON for product", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Set the updated timestamp
	product.ID = uuid.MustParse(productID)
	product.UpdatedAt = time.Now()

	// Perform the DB operation to update the product
	if err := api.Q.UpdateProduct(c, &product); err != nil {
		log.Error("Error updating product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Product updated successfully", zap.String("product_id", product.ID.String()))
	c.JSON(http.StatusOK, gin.H{"error": false, "product": product})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Param id path string true "Product ID"
// @Security CookieAuth
// @Success 200 {object} gin.H{"error": false, "msg": "product deleted successfully"}
// @Failure 401 {object} gin.H{"error": true, "msg": "unauthorized"}
// @Failure 403 {object} gin.H{"error": true, "msg": "permission denied"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router api/products/{id} [delete]
func (api *API) DeleteProduct(c *gin.Context) {
	log := logger.Get()

	// Extract claims from the token
	claims, err := auth.ExtractTokenMetadata(c, "access")
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Check token expiration
	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	// Check if the user has permission to delete the product
	if !claims.Credentials["product:delete"] {
		log.Warn("Permission denied for product deletion", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	// Get the product ID from the URL parameter
	productID := c.Param("id")

	// Perform the DB operation to delete the product
	if err := api.Q.DeleteProduct(c, uuid.MustParse(productID)); err != nil {
		log.Error("Error deleting product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	log.Info("Product deleted successfully", zap.String("product_id", productID))
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "product deleted successfully"})
}

// ListProducts godoc
// @Summary List all products
// @Description Retrieve a list of all products from the database
// @Tags products
// @Security CookieAuth
// @Success 200 {object} gin.H{"error": false, "products": []query.Product}
// @Failure 401 {object} gin.H{"error": true, "msg": "unauthorized"}
// @Failure 403 {object} gin.H{"error": true, "msg": "permission denied"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router /api/products [get]
func (api *API) ListProducts(c *gin.Context) {
	log := logger.Get()

	// Extract claims from the token
	claims, err := auth.ExtractTokenMetadata(c, "access")
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Check token expiration
	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	// Check if the user has permission to view products
	if !claims.Credentials["product:read"] {
		log.Warn("Permission denied for product listing", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	// Retrieve all products from the DB
	products, err := api.Q.GetAllProducts(c)
	if err != nil {
		log.Error("Error retrieving products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if len(products) == 0 {
		log.Warn("No products found")
		c.JSON(http.StatusOK, gin.H{"error": false, "msg": "no products found"})
		return
	}

	log.Info("Products retrieved successfully", zap.Int("product_count", len(products)))
	c.JSON(http.StatusOK, gin.H{"error": false, "products": products})
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Retrieve a product by its ID from the database
// @Tags products
// @Param id path string true "Product ID"
// @Security CookieAuth
// @Success 200 {object} query.Product
// @Failure 401 {object} gin.H{"error": true, "msg": "unauthorized"}
// @Failure 403 {object} gin.H{"error": true, "msg": "permission denied"}
// @Failure 404 {object} gin.H{"error": true, "msg": "product not found"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router /api/products/{id} [get]
func (api *API) GetProduct(c *gin.Context) {
	log := logger.Get()

	// Extract claims from the token
	claims, err := auth.ExtractTokenMetadata(c, "access")
	if err != nil {
		log.Error("Error extracting token metadata", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	// Check token expiration
	if time.Now().Unix() > claims.Exp {
		log.Warn("Token expired", zap.Int64("expiration", claims.Exp))
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "msg": "unauthorized, token expired"})
		return
	}

	// Check if the user has permission to view the product
	if !claims.Credentials["product:read"] {
		log.Warn("Permission denied for product read", zap.String("role", claims.Role))
		c.JSON(http.StatusForbidden, gin.H{"error": true, "msg": "permission denied"})
		return
	}

	// Get the product ID from the URL parameter
	productID := c.Param("id")

	// Retrieve the product from the DB
	product, err := api.Q.GetProductByID(c, uuid.MustParse(productID))
	if err != nil {
		log.Error("Error retrieving product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	if product == nil {
		log.Warn("Product not found", zap.String("product_id", productID))
		c.JSON(http.StatusNotFound, gin.H{"error": true, "msg": "product not found"})
		return
	}

	log.Info("Product retrieved successfully", zap.String("product_id", productID))
	c.JSON(http.StatusOK, gin.H{"error": false, "product": product})
}
