package main

import (
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/server"
	"go.uber.org/zap"
)

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
