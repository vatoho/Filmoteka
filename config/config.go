package config

import (
	"github.com/joho/godotenv"
)

func GetConfig() error {
	envFilePath := ".env"
	return godotenv.Load(envFilePath)
}
