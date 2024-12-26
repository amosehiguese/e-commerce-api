package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Token struct {
	Access  string
	Refresh string
}

func GenerateNewToken(id string, role Role) (*Token, error) {
	accessToken, err := GenerateAccessToken(id, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func GenerateAccessToken(id string, role Role) (string, error) {
	config := config.Get().JWT
	credentials, err := GetRoleCredentials(role)
	if err != nil {
		return "", err
	}

	minCount, err := strconv.Atoi(config.JwtSecretKeyExp)
	if err != nil {
		return "", err
	}

	claims := make(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minCount)).Unix()
	claims["role"] = role

	claims["product:create"] = false
	claims["product:read"] = false
	claims["product:update"] = false
	claims["product:delete"] = false

	claims["order:create"] = false
	claims["order:read"] = false
	claims["order:update"] = false
	claims["order:cancel"] = false

	for _, credential := range credentials {
		claims[credential] = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GenerateRefreshToken() (string, error) {
	config := config.Get().JWT
	hash := sha256.New()
	refresh := config.JwtRefreshKey + time.Now().String()

	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", err
	}

	hoursCount, err := strconv.Atoi(config.JwtRefreshKeyExp)
	if err != nil {
		return "", err
	}

	expTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expTime

	return t, nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}

func AttachToCookie(c *gin.Context, accessToken string, refreshToken string) {
	cfg := config.Get()
	log := logger.Get()

	accessExp, err := strconv.Atoi(cfg.JWT.JwtSecretKeyExp)
	if err != nil {
		log.Error("Invalid token expiration value", zap.Error(err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	accessExpiresAt := time.Now().Add(time.Duration(accessExp) * time.Hour)

	c.SetCookie(
		"access",
		accessToken,
		int(time.Until(accessExpiresAt).Seconds()),
		"/",
		cfg.Domain,
		cfg.Env == "prod",
		true,
	)

	refreshExpiresAt := time.Now().Add(30 * 24 * time.Hour)

	c.SetCookie(
		"refresh",
		refreshToken,
		int(time.Until(refreshExpiresAt).Seconds()),
		"/",
		cfg.Domain,
		cfg.Env == "prod",
		true,
	)

	log.Info("Access and refresh tokens stored in cookies")
}

func InvalidateTokenCookies(c *gin.Context) {
	log := logger.Get()

	c.SetCookie(
		"access",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.SetCookie(
		"refresh",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	log.Info("Access and refresh tokens cookies invalidated!")
}
