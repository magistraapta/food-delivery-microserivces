package common

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port                string `json:"port"`
	EnableGinConsoleLog bool   `json:"enableGinConsoleLog"`
	EnableGinFileLog    bool   `json:"enableGinFileLog"`
	DatabaseName        string `json:"databaseName"`
	JwtSecretKey        string `json:"jwtSecretKey"`
	JwtIssuer           string `json:"jwtIssuer"`
	JwtExpireDuration   string `json:"jwtExpireDuration"`
	DatabaseUrl         string `json:"databaseUrl"`
}

var (
	ConfigData *Config
)

func LoadConfig() error {
	file, err := os.Open("config/config.json")

	if err != nil {
		return err
	}

	config := new(Config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		return err
	}

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	ConfigData = config
	return nil
}
