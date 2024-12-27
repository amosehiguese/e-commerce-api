package main

import (
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/server"
	"go.uber.org/zap"
)

// @title Ecommerce API
// @version 1.0
// @description This is an Ecommerce API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	log := logger.Get()

	log.Info("Starting the eCommerce API server...")
	if err := server.Start(); err != nil {
		log.Fatal("Server failed to start",
			zap.Error(err),
		)
	}
	log.Info("eCommerce API server exited successfully")
}
