package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	sql, err := os.ReadFile("./go/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func FindSearchResults(searchStr string) ([]Result, error) {
	if searchStr == "" {
		return nil, nil
	}

	rows, err := db.Query("SELECT id, name, link FROM results WHERE name LIKE ?", "%"+searchStr+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.ID, &r.Name, &r.Link); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

type Result struct {
	ID   int
	Name string
	Link string
}
