package handlers

import (
	"DVK-Project/db"
	"html/template"
	"net/http"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("go/templates/login.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	ok, err := db.CheckCredentials(email, password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("login successful"))
}
