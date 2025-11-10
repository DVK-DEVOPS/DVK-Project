package db

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")

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
