package main

import (
	"DVK-Project/db"
	"DVK-Project/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)

	r.HandleFunc("/login", handlers.ShowLogin).Methods("GET")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	rc := &handlers.RegistrationController{}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")

	sc := &handlers.SearchController{}
	//r.HandleFunc("/search", sc.ShowSearchPage).Methods("GET")
	r.HandleFunc("/api/search", sc.ShowSearchResults).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
