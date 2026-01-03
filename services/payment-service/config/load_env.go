package config

import "github.com/joho/godotenv"

func LoadEnv() error {
	err := godotenv.Load(".env.local.docker")

	if err != nil {
		return err
	}

	return nil
}
