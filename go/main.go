package main

import (
	"DVK-Project/db"
	"DVK-Project/handlers"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	database := db.Connect()
	userRepository := &db.UserRepository{DB: database}
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	rc := &handlers.RegistrationController{UserRepository: userRepository}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Home Page",
		"Body":  "Hello from Gorilla Mux with HTML templates!",
	}
	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
