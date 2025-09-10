package db

import (
	"database/sql"
	"errors"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) CheckCredentialsByEmail(email, password string) (bool, error) {
	var storedPassword string
	err := ar.db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return storedPassword == HashPassword(password), nil
}
