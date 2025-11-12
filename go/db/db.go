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

	log.Println("[InitDB] Loading environment variables...")
	_ = godotenv.Load()

	log.Println("[InitDB] Getting database URL secret...")
	dbURL := config.GetSecret("DB_URL", os.Getenv("AZURE_KEYVAULT_NAME"), os.Getenv("DB_URL_SECRET_NAME"))

	log.Println("[InitDB] Opening database connection...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("[InitDB] Error opening database connection: %v\n", err)
		return nil, err
	}

	log.Println("[InitDB] Database connection established successfully.")
	return db, nil
}
