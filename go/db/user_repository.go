package db

import (
	"DVK-Project/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CheckIfUserExists(email string) (bool, error) {
	return false, nil
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
