package config

import "github.com/joho/godotenv"

func LoadEnv() error {
	// Optional: when running in K8s, env is set by Deployment; .env is for local/docker-compose
	_ = godotenv.Load(".env.local.docker")
	return nil
}
