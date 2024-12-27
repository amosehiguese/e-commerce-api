package routes

import (
	"database/sql"

	"github.com/amosehiguese/ecommerce-api/api"
	"github.com/amosehiguese/ecommerce-api/middleware"
	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/query"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/amosehiguese/ecommerce-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUp(dbconn *sql.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware
	router.Use(cors.Default())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Initialize Query
	q := query.NewQuery(dbconn)

	// Initialize API
	a := api.NewAPI(q, cfg)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health
	router.GET("/_healthz", a.HealthCheck)

	// Public routes
	public := router.Group("/api/auth")
	RegisterAuthRoutes(public, a)

	// Protected routes (authentication required)
	auth := router.Group("/api", middleware.JWTProtected())
	{
		RegisterTokenRenewal(auth, a)
		RegisterProductRoutes(auth, a)
		RegisterOrderRoutes(auth, a)
	}

	return router
}
