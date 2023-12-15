package config

import (
	"github.com/joho/godotenv"
	"log"
)

func SetEnv() bool {
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatalf("Config-Exception: %s", err)
		return false
	}
	return true
}
