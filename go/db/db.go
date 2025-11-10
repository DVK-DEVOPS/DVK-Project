package db

import (
	"DVK-Project/config"
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	_ = godotenv.Load()
	dbURL := config.GetSecret("DB_URL", os.Getenv("KEYVAULT_NAME"), os.Getenv("DB_URL_SECRET_NAME"))

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	sql, err := os.ReadFile("/whoknows/go/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		return nil, err
	}

	return db, nil
}
