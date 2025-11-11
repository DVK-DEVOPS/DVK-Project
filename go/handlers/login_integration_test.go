package handlers

import (
	"DVK-Project/db"
	"DVK-Project/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func setupTestDB() *db.UserRepository {
	sqlDB, _ := sql.Open("sqlite", ":memory:")
	sqlDB.Exec(`CREATE TABLE users (username TEXT, password TEXT);`)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	sqlDB.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, "test", string(hashed))

	return db.NewUserRepository(sqlDB)
}

func TestLoginIntegrationSuccess(t *testing.T) {
	repo := setupTestDB()
	lh := &LoginHandler{UserRepository: repo}

	data := "username=test&password=pass"
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	lh.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp models.AuthResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.StatusCode != 3070 {
		t.Errorf("expected StatusCode 3070, got %d", resp.StatusCode)
	}
}

func TestLoginIntegrationFail(t *testing.T) {
	repo := setupTestDB()
	lh := &LoginHandler{UserRepository: repo}

	data := "username=test&password=wrong"
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	lh.Login(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}
