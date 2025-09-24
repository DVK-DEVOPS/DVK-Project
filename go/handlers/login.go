package handlers

import (
	"DVK-Project/db"
	"DVK-Project/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

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

// Login authenticates a user with username and password.
// @Summary Login
// @Description Authenticates a user using username and password
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} models.AuthResponse "Successful registration"
// @Failure 422 {object} models.HTTPValidationError "Validation error"
// @Router /api/login [post]
func (lh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var username, password string

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(models.HTTPValidationError{
				Detail: []models.ValidationErrorDetail{
					{Loc: []interface{}{"body"}, Msg: "Invalid JSON body", Type: "parse_error"},
				},
			})
			return
		}
		username, password = creds.Username, creds.Password
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(models.HTTPValidationError{
				Detail: []models.ValidationErrorDetail{
					{Loc: []interface{}{"body", "form"}, Msg: "Form parse error", Type: "parse_error"},
				},
			})
			return
		}
		username = r.FormValue("username")
		password = r.FormValue("password")
	}

	if username == "" || password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		details := []models.ValidationErrorDetail{}
		if username == "" {
			details = append(details, models.ValidationErrorDetail{
				Loc:  []interface{}{"body", "username"},
				Msg:  "Username required",
				Type: "validation_error",
			})
		}
		if password == "" {
			details = append(details, models.ValidationErrorDetail{
				Loc:  []interface{}{"body", "password"},
				Msg:  "Password required",
				Type: "validation_error",
			})
		}
		json.NewEncoder(w).Encode(models.HTTPValidationError{Detail: details})
		return
	}

	ok, err := lh.UserRepository.CheckCredentialsByUsername(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{Loc: []interface{}{"db"}, Msg: "Database error", Type: "internal_error"},
			},
		})
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{Loc: []interface{}{"body", "credentials"}, Msg: "Invalid credentials", Type: "validation_error"},
			},
		})
		return
	}

	value := map[string]string{"username": username}
	encoded, _ := s.Encode("session", value)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		StatusCode: 3070,
		Message:    fmt.Sprintf("User authenticated with username %s", username),
	})
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
