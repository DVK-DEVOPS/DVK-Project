package handlers

import (
	"DVK-Project/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockUserRepo struct{}

func (m *mockUserRepo) CheckCredentialsByUsername(username, password string) (bool, error) {
	return username == "test" && password == "pass", nil
}

type LoginHandlerUnit struct {
	UserRepository *mockUserRepo
}

func (lh *LoginHandlerUnit) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ok, _ := lh.UserRepository.CheckCredentialsByUsername(username, password)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(models.AuthResponse{
		StatusCode: 3070,
		Message:    "User authenticated",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func TestLoginUnit_Success(t *testing.T) {
	handler := &LoginHandlerUnit{UserRepository: &mockUserRepo{}}

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader("username=test&password=pass"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp models.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if resp.StatusCode != 3070 {
		t.Errorf("expected 3070, got %d", resp.StatusCode)
	}
}

func TestLoginUnit_InvalidCredentials(t *testing.T) {
	handler := &LoginHandlerUnit{UserRepository: &mockUserRepo{}}

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader("username=test&password=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestLoginUnit_MissingFields(t *testing.T) {
	handler := &LoginHandlerUnit{UserRepository: &mockUserRepo{}}

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader("username="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}
