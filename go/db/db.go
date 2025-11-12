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
	_ = godotenv.Load() // loads .env if present

	dbURL := ""
	if os.Getenv("APP_ENV") == "dev" {
		log.Println("[InitDB] Using local dev DB URL...")
		dbURL = os.Getenv("DB_URL")
	} else {
		log.Println("[InitDB] Using production DB URL from Key Vault...")
		dbURL = config.GetSecret(
			"DB_URL",
			os.Getenv("AZURE_KEYVAULT_NAME"),
			os.Getenv("DB_URL_SECRET_NAME"),
		)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("[InitDB] Database connection established.")
	return db, nil
}
