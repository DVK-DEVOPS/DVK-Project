package main

import (
	handlers2 "DVK-Project/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/register", handlers2.RegistrationController{}.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", handlers2.RegistrationController{}.Register).Methods("POST")
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
