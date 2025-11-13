package db

import (
	"DVK-Project/config"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	var dbURL string

	if env == "production" {
		log.Println("[InitDB] Loading production DB URL from Key Vault...")
		dbURL = config.GetSecret(
			"DB_URL",
			os.Getenv("AZURE_KEYVAULT_NAME"),
			os.Getenv("DB_URL_SECRET_NAME"),
		)
	} else {
		log.Println("[InitDB] Using local development DB URL from .env...")
		dbURL = os.Getenv("DB_URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("[InitDB] Database connection established successfully.")
	return db, nil
}
