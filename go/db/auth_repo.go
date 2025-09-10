package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
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
	return storedPassword == hashPassword(password), nil
}
