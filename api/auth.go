package api

import (
	"context"
	"net/http"
	"time"

	"github.com/amosehiguese/ecommerce-api/api/payload"
	"github.com/amosehiguese/ecommerce-api/pkg/auth"
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/pkg/utils"
	"github.com/amosehiguese/ecommerce-api/pkg/validator"
	"github.com/amosehiguese/ecommerce-api/query"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (a *API) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var registerPayload payload.RegisterPayload
	if err := c.ShouldBindJSON(&registerPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}
	registerPayload.Role = "user"

	validate := validator.NewValidator()
	if err := validate.Struct(registerPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	role, err := auth.VerifyRole(registerPayload.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	user := &query.User{
		ID:           uuid.New(),
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
		FirstName:    registerPayload.FirstName,
		LastName:     &registerPayload.LastName,
		Email:        registerPayload.Email,
		PasswordHash: utils.HashPassword(registerPayload.Password),
		Role:         role.String(),
	}

	validate = validator.NewValidator()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	user, err = a.Q.CreateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Unable to create user: " + err.Error(),
		})
		return
	}

	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "User created successfully",
		"user":  user.ID,
	})
}

func (api *API) Login(c *gin.Context) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var loginPayload payload.LoginPayload
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	validate := validator.NewValidator()
	if err := validate.Struct(loginPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	// Get user by email
	u, err := api.Q.GetUserByEmail(ctx, loginPayload.Email)
	if err != nil {
		log.Error("user with the given email not found", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "User with the given email is not found",
		})
		return
	}

	role, err := auth.VerifyRole(u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
		})
		return
	}

	_, err = auth.GenerateNewToken(u.ID.String(), role)
	if err != nil {
		log.Error("Failed to generate tokens", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to generate tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": "JWT_TOKEN_HERE"})
}
