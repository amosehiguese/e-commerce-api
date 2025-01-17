package config

import (
	"strings"

	"github.com/amosehiguese/ecommerce-api/pkg/utils"
)

type jwtConfig struct {
	JwtSecretKey     string
	JwtSecretKeyExp  string
	JwtRefreshKey    string
	JwtRefreshKeyExp string
	CorsOrigins      []string
}

func setJwtConfig() *jwtConfig {
	var s jwtConfig
	utils.MustMapEnv(&s.JwtSecretKey, "JWT_SECRET_KEY")
	utils.MustMapEnv(&s.JwtSecretKeyExp, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT")
	utils.MustMapEnv(&s.JwtRefreshKey, "JWT_REFRESH_KEY")
	utils.MustMapEnv(&s.JwtRefreshKeyExp, "JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT")

	var coreStr string
	utils.MustMapEnv(&coreStr, "CORS_ORIGINS")
	if coreStr != "" {
		s.CorsOrigins = strings.Split(coreStr, ",")
	}

	return &s
}
