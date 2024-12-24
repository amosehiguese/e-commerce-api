package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	env := os.Getenv("ECOMM_ENV")
	if env == "" {
		env = "development"
	}

	switch env {
	case "production":
		err := godotenv.Load(".env.prod")
		if err != nil {
			log.Printf("No .env.prod file found")
		}

	case "development":
		err := godotenv.Load(".env.dev")
		if err != nil {
			log.Printf("No .env.dev file")
		}
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func MustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}

	*target = v
}
