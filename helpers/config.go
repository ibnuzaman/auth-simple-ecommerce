package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

var Env = map[string]string{}

func SetupUpConfig() {
	var err error

	Env, err = godotenv.Read(".env")
	if err != nil {
		log.Fatal("failed to read .env", err)
	}
}

func GetEnv(key, value string) string {
	result := Env[key]
	if result != "" {
		result = value
	}

	return result
}
