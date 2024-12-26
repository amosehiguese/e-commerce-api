package auth

import (
	"errors"

	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenMetadata struct {
	UserID      uuid.UUID
	Credentials map[string]bool
	Role        string
	Exp         int64
}

func ExtractTokenMetadata(c *gin.Context, name string) (*TokenMetadata, error) {
	token, err := verifyToken(c, name)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			return nil, err
		}

		exp := int64(claims["exp"].(float64))
		role := claims["role"].(string)

		credentials := map[string]bool{
			// Product Permissions
			"product:create": claims["product:create"].(bool),
			"product:read":   claims["product:read"].(bool),
			"product:update": claims["product:update"].(bool),
			"product:delete": claims["product:delete"].(bool),

			// Order Permissions
			"order:create": claims["order:create"].(bool),
			"order:read":   claims["order:create"].(bool),
			"order:update": claims["order:create"].(bool),
			"order:cancel": claims["order:create"].(bool),
		}

		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Role:        role,
			Exp:         exp,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func verifyToken(c *gin.Context, name string) (*jwt.Token, error) {
	config := config.Get().JWT
	tokenString, err := extractToken(c, name)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func extractToken(c *gin.Context, name string) (string, error) {
	token, err := c.Cookie(name)
	if err != nil {
		return "", errors.New("missing JWT")
	}

	return token, nil
}
