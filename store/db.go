package store

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/pkg/logger"
)

func SetUpDB(config *config.Config) (*sql.DB, error) {
	dbCfg := config.Database
	connStr := dbCfg.ConnStringDefaultDB()

	log := logger.Get()
	log.Info("Connecting to database server",
		zap.String("host", dbCfg.Host),
		zap.Int("port", dbCfg.Port),
	)

	// Connect to the server (without specifying the database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}
	defer db.Close()

	// Check server connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Postgres server: %w", err)
	}

	log.Info("Successfully connected to Postgres server", zap.String("user", dbCfg.User))

	// Create database if it doesn't exist
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbCfg.Name)
	var exists bool
	if err := db.QueryRow(query).Scan(&exists); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Info("Database does not exist, creating database", zap.String("database_name", dbCfg.Name))
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbCfg.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		log.Info("Database created successfully", zap.String("database_name", dbCfg.Name))
	} else {
		log.Info("Database already exists", zap.String("database_name", dbCfg.Name))
	}

	// Connect to the newly created or existing database
	connStr = dbCfg.ConnString()
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres database: %w", err)
	}

	log.Info("Connected to target database", zap.String("database_name", dbCfg.Name))

	// Apply migrations
	err = applyMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Info("Migrations applied successfully", zap.String("database_name", dbCfg.Name))
	return db, nil
}

func applyMigrations(db *sql.DB) error {
	projectRoot, _ := findProjectRoot()
	migrationsPath, err := findMigrationsDir(projectRoot)
	if err != nil {
		return err
	}

	err = goose.Up(db, migrationsPath)
	if err != nil {
		return err
	}
	return nil
}

func findMigrationsDir(startingDir string) (string, error) {
	migrationsDir := filepath.Join(startingDir, "store/migrations")

	if stat, err := os.Stat(migrationsDir); err == nil && stat.IsDir() {
		return migrationsDir, nil
	}

	return "", os.ErrNotExist
}

func findProjectRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Git project root: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}
