package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./DVK-project.db")
	if err != nil {
		return err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT
	); 
	CREATE TABLE IF NOT EXISTS results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, 
		link TEXT)
		`)
	return err
}

func CheckCredentials(username, password string) (bool, error) {
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return storedPassword == password, nil
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
