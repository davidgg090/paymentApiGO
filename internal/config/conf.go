package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DSN    string
	JWTKey []byte
)

func LoadConfig() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	DSN = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_PASSWORD"),
	)

	JWTKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}
