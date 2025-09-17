package handlers

import (
	"DVK-Project/db"
	"html/template"
	"net/http"
)

type LoginHandler struct {
	UserRepository *db.UserRepository
}

// ShowLogin
// @Summary Show login page to the user
// @Description showing the login page to the user
// @Tags users
// @Router /login [post]

func (lh *LoginHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/templates/login.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Login authenticates a user with email and password.
// @Summary Login
// @Description Authenticates a user using email and password
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce text/plain
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 200 {object} models.AuthResponse "Successful registration"
// @Failure 422 {object} models.HTTPValidationError "Validation error"
// @Router /api/login [post]
func (lh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	ok, err := lh.UserRepository.CheckCredentialsByEmail(email, password)
	if err != nil {
		http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("login successful"))
}
