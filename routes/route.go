package routes

import (
	"database/sql"

	"github.com/amosehiguese/ecommerce-api/api"
	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUp(dbconn *sql.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware
	router.Use(cors.Default())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Routes
	router.GET("/_healthz", api.HealthCheck)

	return router
}
