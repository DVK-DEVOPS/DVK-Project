package handlers

import (
	"DVK-Project/db"
	"DVK-Project/models"
	"encoding/json"
	"html/template"
	"net/http"
)

type RegistrationController struct {
	UserRepository *db.UserRepository
}

func (rc *RegistrationController) ShowRegistrationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("go/templates/register.html")
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
// @Success 200 {string} string "User created with ID"
// @Failure 422 {string} string "Validation error"
// @Router /api/register [post]
func (rc *RegistrationController) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Form parse error",
		})
		if err != nil {
			return
		}
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if username == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "All fields are required",
		})
		return
	}

	hashedPassword, err := db.HashPassword(password)

	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	exists, err := rc.UserRepository.CheckIfUserExists(user.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	if exists {
		http.Error(w, "User with this email is already registered", http.StatusUnprocessableEntity)
	}

	id, err := rc.UserRepository.AddUser(user)
	if err != nil {
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"id":      id,
	})

}
