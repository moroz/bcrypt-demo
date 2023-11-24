package config

import (
	"log"
	"os"
)

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("FATAL: Environment variable %s is not set!\n", key)
	}
	return value
}

var DATABASE_URL = MustGetEnv("DATABASE_URL")
