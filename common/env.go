package common

import (
	"log"

	"github.com/joho/godotenv"
)

var appEnv map[string]string

func LoadENV() {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading env")
	}

	if env["MONGO_URI"] == "" {
		log.Fatal("No  MONGO_URI env Variable found")
	}

	if env["MONGO_DATABASE"] == "" {
		log.Fatal("No MONGO_DATABASE env Variable found")
	}

	if env["JWT_SECRET"] == "" {
		log.Fatal("No JWT_SECRET env Variable found")
	}

	if env["S3_BUCKET"] == "" {
		log.Fatal("No S3_BUCKET env Variable found")
	}

	if env["S3_REGION"] == "" {
		log.Fatal("No S3_REGION env Variable found")
	}

	appEnv = env
}

func GetENV(key string) string {
	return appEnv[key]
}