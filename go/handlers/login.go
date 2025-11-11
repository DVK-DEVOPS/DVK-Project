package handlers

import (
	"DVK-Project/db"
	"DVK-Project/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
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
	renderTemplate(w, r, "login.html", nil)
}

func (lh *LoginHandler) ShowPasswordReset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Accessing ShowPasswordReset() Showing password-reset.html")
	renderTemplate(w, r, "password-reset.html", nil)
}

func (lh *LoginHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var username, email, oldPassword, newPasswordStr, newPassword string

	//_ = r.ParseForm()
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	username = r.Form.Get("username")
	//email = r.Form.Get("email")
	oldPassword = r.Form.Get("password")
	newPasswordStr = r.Form.Get("newPassword")

	username = strings.TrimSpace(username)
	//email = strings.TrimSpace(email)
	oldPassword = strings.TrimSpace(oldPassword)
	newPasswordStr = strings.TrimSpace(newPassword)
	newPassword, _ = db.HashPassword(newPasswordStr)

	// Prepare data to pass to template
	data := map[string]interface{}{
		"Username": username,
		"Email":    email,
		"Error":    "", // default no error
	}

	ok, err := lh.UserRepository.CheckCredentialsByUsername(username, oldPassword)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !ok {
		data["Error"] = "Old password is incorrect"
		renderTemplate(w, r, "password-reset.html", data)
		return
	}
	// Update password
	//err = lh.UserRepository.UpdatePassword(username, newPassword)
	var rowsAffected int64
	rowsAffected, err = lh.UserRepository.UserResetPassword(username, newPassword)
	if err != nil {
		data["Error"] = "Failed to update password"
		renderTemplate(w, r, "password-reset.html", data)
		return
	}

	if rowsAffected == 0 {
		data["Error"] = "No user found with that username"
		renderTemplate(w, r, "password-reset.html", data)
		return
	}

	data["Success"] = "Password updated successfully"
	renderTemplate(w, r, "password-reset.html", data)
}

// Login authenticates a user with username and password.
// @Summary Login
// @Description Authenticates a user using username and password
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {object} models.AuthResponse "Successful registration"
// @Failure 422 {object} models.HTTPValidationError "Validation error"
// @Router /api/login [post]
func (lh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	isBrowser := strings.Contains(r.Header.Get("Accept"), "text/html")

	w.Header().Set("Content-Type", "application/json")

	var username, password string

	_ = r.ParseForm()
	username = r.Form.Get("username")
	password = r.Form.Get("password")

	// If form fields not present, try JSON
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err == nil {
			if username == "" {
				username = creds.Username
			}
			if password == "" {
				password = creds.Password
			}
		}
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if lh.UserRepository.CheckIfUserIsAffected(username) { //Is user affected by security breach?
		fmt.Println("(Login.go)Checking if user is affected by security breach.")
		if isBrowser {
			fmt.Println("(Login.go) !AFFECTED! User is accessing through browser. Trying redirect.")
			http.Redirect(w, r, "/password-reset", http.StatusFound)
			return
		}
		fmt.Println("(Login.go) !AFFECTED! User is accessing via api and json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Password reset required",
			"url":   "/password-reset",
		})
		return
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
		_ = json.NewEncoder(w).Encode(models.HTTPValidationError{Detail: details})
		return
	}

	ok, err := lh.UserRepository.CheckCredentialsByUsername(username, password)
	if err != nil {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{Loc: []interface{}{"db"}, Msg: "Database error", Type: "internal_error"},
			},
		})
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(models.HTTPValidationError{
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
	_ = json.NewEncoder(w).Encode(models.AuthResponse{
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
