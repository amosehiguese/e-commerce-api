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

	auth.AttachToCookie(c, token.Refresh)

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"tokens": gin.H{
			"access":  token.Access,
			"refresh": token.Refresh,
		},
		"user": gin.H{
			"id":    u.ID,
			"email": u.Email,
			"role":  u.Role,
		},
	})
}

// CreateAdmin godoc
// @Summary      Create an admin user
// @Description  Creates a new admin user with a hashed password.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        admin body payload.RegisterPayload true "Admin creation payload"
// @Success      200 {object} map[string]interface{} "Admin user created successfully"
// @Failure      400 {object} map[string]interface{} "Invalid request data"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Security     BearerAuth
// @Router       api/auth/create-admin [post]
func (api *API) CreateAdmin(c *gin.Context) {
	log := logger.Get()

	var adminPayload payload.RegisterPayload
	if err := c.ShouldBindJSON(&adminPayload); err != nil {
		log.Error("Invalid JSON for creating admin", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	adminPayload.Role = "admin"

	validate := validator.NewValidator()
	if err := validate.Struct(adminPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	role, err := auth.VerifyRole(adminPayload.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	admin := &query.User{
		ID:           uuid.New(),
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
		FirstName:    adminPayload.FirstName,
		LastName:     &adminPayload.LastName,
		Email:        adminPayload.Email,
		PasswordHash: utils.HashPassword(adminPayload.Password),
		Role:         role.String(),
	}

	validate = validator.NewValidator()
	if err := validate.Struct(admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   validator.ValidatorErrors(err),
		})
		return
	}

	admin, err = api.Q.CreateUser(c, admin)
	if err != nil {
		log.Error("Error creating admin user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	admin.PasswordHash = ""

	log.Info("Admin user created successfully", zap.String("username", admin.FirstName))
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "Admin user created successfully"})
}

func (api *API) RenewTokens(c *gin.Context) {
	logger := logger.Get()
	now := time.Now().Unix()

	claims, err := auth.ExtractTokenMetadata(c)
	if err != nil {
		logger.Error("Error extracting token metadata", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	expiresAccessToken := claims.Exp

	if now > expiresAccessToken {
		logger.Warn("Access token expired", zap.Int64("expires_at", expiresAccessToken))

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
		return
	}

	renew := &Renew{}

	if err := c.ShouldBindJSON(renew); err != nil {
		logger.Error("Error binding JSON body", zap.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	expiresRefreshToken, err := auth.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		logger.Error("Error parsing refresh token", zap.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"msg":   err.Error(),
		})
		return
	}

	if now < expiresRefreshToken {
		userID := claims.UserID

		u, err := api.Q.GetUserByID(c, userID)
		if err != nil {
			logger.Error("User not found", zap.String("user_id", userID.String()))

			c.JSON(http.StatusNotFound, gin.H{
				"error": true,
				"msg":   "user with the given ID is not found",
			})
			return
		}

		role, err := auth.VerifyRole(u.Role)
		if err != nil {
			logger.Error("Error getting credentials by role", zap.String("role", u.Role), zap.String("error", err.Error()))

			c.JSON(http.StatusBadRequest, gin.H{
				"error": true,
				"msg":   err.Error(),
			})
			return
		}

		tokens, err := auth.GenerateNewToken(userID.String(), role)
		if err != nil {
			logger.Error("Error generating new tokens", zap.String("error", err.Error()))

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"msg":   err.Error(),
			})
			return
		}

		// Create a new Redis connection
		auth.AttachToCookie(c, tokens.Refresh)

		// Log the successful token renewal
		logger.Info("Token renewal successful", zap.String("user_id", userID.String()))

		// Return tokens
		c.JSON(http.StatusOK, gin.H{
			"error": false,
			"msg":   nil,
			"tokens": gin.H{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	} else {
		// Log the session ended error
		logger.Warn("Unauthorized, session ended earlier")

		// Return status 401 and unauthorized error message
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"msg":   "unauthorized, your session was ended earlier",
		})
	}
}

type Renew struct {
	RefreshToken string `json:"refresh_token"`
}
