package handlers

import (
	"html/template"
	"net/http"
)

type RegistrationController struct {
}

func (rc *RegistrationController) ShowRegistrationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/registration.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

func (rc *RegistrationController) Register(w http.ResponseWriter, r *http.Request) {

}
