package db

import (
	"crypto/md5"
	"encoding/hex"
)

// TODO move utils.go to another package for more seperation
func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
