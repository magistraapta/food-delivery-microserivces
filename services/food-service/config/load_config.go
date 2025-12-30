package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}

	return nil
}
