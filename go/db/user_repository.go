package db

import (
	"DVK-Project/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CheckIfUserExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	err := r.DB.QueryRow(query, email).Scan(&count)

	return false, err
}

func (r *UserRepository) AddUser(user models.User) (int, error) {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	res, err := r.DB.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
