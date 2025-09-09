package handlers

import (
	"html/template"
	"net/http"
)

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

}
