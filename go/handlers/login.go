package handlers

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("go/templates/*.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Login",
		"Body":  "Hello from Gorilla Mux with HTML templates!",
	}
	err := tmpl.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
