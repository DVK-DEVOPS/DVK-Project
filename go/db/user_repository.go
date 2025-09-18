package db

import (
	"DVK-Project/models"
	"database/sql"
	"errors"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CheckIfUserExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	err := r.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
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

// for login
func (r *UserRepository) CheckCredentialsByUsername(username, password string) (bool, error) {
	var storedPassword string
	err := r.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return VerifyPassword(storedPassword, password), nil
}
