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

// Checking isAffected column
func (r *UserRepository) CheckIfUserIsAffected(username string) bool {
	var isAffected bool
	err := r.DB.QueryRow("SELECT isAffected FROM users WHERE username = ?", username).Scan(&isAffected) //TODO: awaiting final column name
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false //User not found
		}
		return false //DB error
	}
	return isAffected
}

// Reset password
func (r *UserRepository) UserResetPassword(username, newPassword string) (int64, error) {
	//hashedPassword := db.HashPassword(newPassword)
	query := "UPDATE USERS SET PASSWORD = ? WHERE USERNAME = ?" //Update the boolean column
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
