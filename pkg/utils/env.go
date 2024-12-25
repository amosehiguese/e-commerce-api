package utils

import (
	"fmt"
	"os"
	"strconv"
)

func MustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}

	*target = v
}

func GetEnvAsInt(envKey string) int {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}

	val, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("failed to convert %q", envKey))
	}

	return val
}
