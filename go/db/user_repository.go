package db

import (
	"DVK-Project/models"
	"database/sql"
	"errors"
	//"os/user"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}
func (r *UserRepository) CheckIfUserExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	err := r.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) AddUser(user models.User) (int, error) {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := r.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) CheckCredentialsByUsername(username, password string) (bool, error) {
	var storedPassword string
	err := r.DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return VerifyPassword(storedPassword, password), nil
}

// Checking is_inactive column
func (r *UserRepository) CheckIfUseris_inactive(username string) bool {
	var is_inactive bool
	err := r.DB.QueryRow("SELECT is_inactive FROM users WHERE username = $1", username).Scan(&is_inactive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false //User not found
		}
		return false //DB error
	}
	return is_inactive
}

// Reset password
func (r *UserRepository) UserResetPassword(username, newPassword string) (int64, error) {
	//hashedPassword := db.HashPassword(newPassword)
	query := "UPDATE USERS SET PASSWORD = $1, is_inactive = false WHERE USERNAME = $2"
	res, err := r.DB.Exec(query, newPassword, username)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil //return 1 or 0

}
