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

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param registerPayload body payload.RegisterPayload true "User Registration Data"
// @Success 200 {object} gin.H{"error": false, "msg": "User created successfully", "user": "user_id"}
// @Failure 400 {object} gin.H{"error": true, "msg": "error message"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router /api/auth/register [post]
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

// Login godoc
// @Summary User Login
// @Description Logs in an existing user and returns the access token
// @Tags auth
// @Accept json
// @Produce json
// @Param loginPayload body payload.LoginPayload true "User Login Data"
// @Success 200 {object} gin.H{"error": false, "tokens": {"access": "access_token"}, "user": {"id": "user_id", "email": "email", "role": "role"}}
// @Failure 400 {object} gin.H{"error": true, "msg": "error message"}
// @Failure 500 {object} gin.H{"error": true, "msg": "error message"}
// @Router /api/auth/login [post]
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

	if !u.ComparePasswordHash(loginPayload.Password) {
		log.Error("Password mismatch", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   "Invalid credentials",
		})
		return
	}

	token, err := auth.GenerateNewToken(u.ID.String(), role)
	if err != nil {
		log.Error("Failed to generate tokens", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   "Failed to generate tokens",
		})
		return
	}

	auth.AttachToCookie(c, token.Access, token.Refresh)

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"tokens": gin.H{
			"access": token.Access,
		},
		"user": gin.H{
			"id":    u.ID,
			"email": u.Email,
			"role":  u.Role,
		},
	})
}

// Logout godoc
// @Summary User Logout
// @Description Logs out the current user and invalidates their session
// @Tags auth
// @Security CookieAuth
// @Success 204 {object} gin.H{}
// @Failure 400 {object} gin.H{"error": true, "msg": "error message"}
// @Router api/auth/logout [post]
func (api *API) Logout(c *gin.Context) {
	auth.InvalidateTokenCookies(c)
	c.Status(204)
}
