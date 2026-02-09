package config

import "github.com/joho/godotenv"

func LoadConfig() error {
	// .env is optional when running in K8s/containers (env is injected via ConfigMap/Secret)
	_ = godotenv.Load(".env")
	return nil
}
