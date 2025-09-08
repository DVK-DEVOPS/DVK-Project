package main

import (
	"DVK-Project/handlers"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseGlob("go/templates/*.html"))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/login", handlers.ShowLogin)

	rc := &handlers.RegistrationController{}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
