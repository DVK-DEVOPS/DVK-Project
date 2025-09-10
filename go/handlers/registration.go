package handlers

import (
	"DVK-Project/db"
	"DVK-Project/models"
	"fmt"
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

func (rc *RegistrationController) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	hashedPassword, err := db.HashPassword(r.FormValue("password"))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: hashedPassword,
	}

	exists, err := rc.UserRepository.CheckIfUserExists(user.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

	if exists {
		http.Error(w, "User with this email is already registered", http.StatusBadRequest)
	}

	id, err := rc.UserRepository.AddUser(user)
	if err != nil {
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User created with id: %d", id)
}
