package handlers

import (
	"DVK-Project/db"
	"DVK-Project/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type RegistrationController struct {
	UserRepository *db.UserRepository
}

// @Summary Serve register page
// @Description Show the registration page.
// @Tags users
// @Produce text/html
// @Success 200 {string} string "HTML of registration page"
// @Failure 404 {string} string "Template not found"
// @Router /register [get]
func (rc *RegistrationController) ShowRegistrationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/register.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

// @Summary Register new user
// @Description Create new user with the parameters user fills in the registration form.
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {object} models.AuthResponse "Successful registration"
// @Failure 422 {object} models.HTTPValidationError "Validation error"
// @Router /api/register [post]
func (rc *RegistrationController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{
					Loc:  []interface{}{"body", "form"},
					Msg:  "Form parse error",
					Type: "parse_error",
				},
			},
		})
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{
					Loc:  []interface{}{"body", "fields"},
					Msg:  "All fields are required",
					Type: "validation_error",
				},
			},
		})
		return
	}

	hashedPassword, err := db.HashPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{Loc: []interface{}{"db"}, Msg: "Internal server error", Type: "internal_error"},
			},
		})
		return
	}
	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	exists, err := rc.UserRepository.CheckIfUserExists(user.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{
					Loc:  []interface{}{"db"},
					Msg:  "Database error",
					Type: "internal_error"},
			},
		})
	}

	if exists {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{
					Loc:  []interface{}{"body", "email"},
					Msg:  "User with this email already exists",
					Type: "validation_error",
				},
			},
		})
	}

	id, err := rc.UserRepository.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.HTTPValidationError{
			Detail: []models.ValidationErrorDetail{
				{
					Loc:  []interface{}{"db"},
					Msg:  "Error adding user",
					Type: "internal_error"},
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("User created with ID %d", id),
	})

}
