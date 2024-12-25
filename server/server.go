package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
	"github.com/amosehiguese/ecommerce-api/pkg/utils"
	"github.com/amosehiguese/ecommerce-api/routes"
	"github.com/amosehiguese/ecommerce-api/store"
	"go.uber.org/zap"
)

func Start() error {
	cfg := config.Get()
	log := logger.Get()

	// Prepare listener address
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Info("Preparing to bind to address", zap.String("address", addr))

	// Bind to address
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to bind to address",
			zap.String("address", addr),
			zap.Error(err),
		)
		return err
	}
	defer listener.Close()
	log.Info("Successfully bound to address", zap.String("address", addr))

	// Set up database connection
	log.Info("Configuring database", zap.String("database_name", cfg.Database.Name))
	dbConn, err := store.SetUpDB(cfg)
	if err != nil {
		log.Error("Failed to configure database",
			zap.String("database_name", cfg.Database.Name),
			zap.Error(err),
		)
		return err
	}
	defer dbConn.Close() // Ensure database connection is closed
	log.Info("Database configured successfully", zap.String("database_name", cfg.Database.Name))

	// Run application
	log.Info("Starting application...", zap.String("address", addr), zap.String("environment", cfg.Env))
	err = run(listener, dbConn, cfg)
	if err != nil {
		log.Error("Application failed to start", zap.Error(err))
		return err
	}

	return nil
}

func run(l net.Listener, dbconn *sql.DB, cfg *config.Config) error {
	// SetUp Router
	router := routes.SetUp(dbconn, cfg)
	server := &http.Server{
		Addr:    l.Addr().String(),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			logger.Get().Error("Server failed",
				zap.String("address", server.Addr),
				zap.Error(err),
			)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit // Wait for interrupt signal

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		logger.Get().Error("Server shutdown failed", zap.Error(err))
		return err
	}

	logger.Get().Info("Server shut down gracefully")
	return nil
}

// For Testing Purposes
type TestApp struct {
	Addr    int
	DB      *sql.DB
	DB_Name string
}

func SpawnApp() (*TestApp, error) {
	log := logger.Get()

	// Bind to a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Error("Failed to bind to a random port", zap.Error(err))
		return nil, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	log.Info("Successfully bound to a random port", zap.Int("port", port))

	// Load configuration
	cfg := config.Get()
	rawUUID := utils.GenUUID()
	cfg.Database.Name = fmt.Sprintf("testdb_%s", rawUUID[:8])
	log.Info("Generated random database name for test", zap.String("database_name", cfg.Database.Name))

	// Set up database
	dbConn, err := store.SetUpDB(cfg)
	if err != nil {
		log.Error("Failed to configure database", zap.String("database_name", cfg.Database.Name), zap.Error(err))
		return nil, err
	}
	log.Info("Database configured successfully", zap.String("database_name", cfg.Database.Name))

	// Run application
	log.Info("Starting application...", zap.String("address", listener.Addr().String()), zap.String("environment", cfg.Env))
	go func() {
		if err := run(listener, dbConn, cfg); err != nil {
			log.Error("Failed to run app", zap.Error(err))
		}
	}()

	log.Info("Application started successfully", zap.String("address", listener.Addr().String()))

	return &TestApp{
		Addr:    port,
		DB:      dbConn,
		DB_Name: cfg.Database.Name,
	}, nil
}

func DropTestDatabase(dbConn *sql.DB, dbName string) error {
	log := logger.Get()
	cfg := config.Get().Database

	defaultDBConn, err := sql.Open("postgres", cfg.ConnStringDefaultDB())
	if err != nil {
		log.Error("Error connecting to the default database: %v", zap.Error(err))
		return err
	}
	defer defaultDBConn.Close()

	// Terminate active connections to the test database before dropping it
	_, err = defaultDBConn.Exec(fmt.Sprintf(`SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid();`, dbName))
	if err != nil {
		log.Error("Error terminating connections: %v", zap.Error(err))
		return err
	}

	// Drop the database
	_, err = defaultDBConn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
	if err != nil {
		log.Error("Error dropping test database: %v", zap.Error(err))
		return err
	}

	log.Sugar().Infof("Test database %s dropped successfully.", dbName)
	return nil
}
