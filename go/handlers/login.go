package handlers

import (
	"DVK-Project/db"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
)

var hashKey = []byte("secure-key")
var s = securecookie.New(hashKey, nil)

type LoginHandler struct {
	UserRepository *db.UserRepository
}

// ShowLogin
// @Summary Show login page
// @Description Displays the login page
// @Tags users
// @Produce text/html
// @Success 200 {string} string "Successful"
// @Failure 404 {string} string "Error"
// @Router /login [get]
func (lh *LoginHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
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
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce application/json
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

	value := map[string]string{"email": email}
	encoded, _ := s.Encode("session", value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout logs the user out by clearing the session cookie.
// @Summary Logout
// @Description Logs the user out by deleting the session cookie and redirects to login page
// @Tags users
// @Produce text/html
// @Success 200 {object} models.AuthResponse "Successful logout"
// @Router /logout [get]
func (lh *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
