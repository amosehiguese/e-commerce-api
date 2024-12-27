package auth

import (
	"errors"
	"strings"

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

func ExtractTokenMetadata(c *gin.Context) (*TokenMetadata, error) {
	token, err := verifyToken(c)
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
			"order:read":   claims["order:read"].(bool),
			"order:update": claims["order:update"].(bool),
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

func verifyToken(c *gin.Context) (*jwt.Token, error) {
	config := config.Get().JWT
	tokenString := extractToken(c)

	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecretKey), nil
	})

	return token, nil
}

func extractToken(c *gin.Context) string {
	bearer := c.GetHeader("Authorization")
	token := strings.Split(bearer, " ")
	if len(token) == 2 {
		return token[1]
	}

	return ""
}
