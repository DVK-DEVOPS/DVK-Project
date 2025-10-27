package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html") //Removing go/ from path serves the html correctly and does not crash the app when navigating to /
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/about.html") //Removing go/ from path serves the html correctly and does not crash the app when navigating to /
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}
