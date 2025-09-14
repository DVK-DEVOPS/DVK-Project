package main

import (
	"DVK-Project/db"
	"DVK-Project/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	userRepository := db.NewUserRepository(database)
	pageRepository := db.NewPageRepository(database)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)

	lh := &handlers.LoginHandler{UserRepository: userRepository}
	r.HandleFunc("/login", lh.ShowLogin).Methods("GET")
	r.HandleFunc("/api/login", lh.Login).Methods("POST")

	rc := &handlers.RegistrationController{UserRepository: userRepository}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")

	sc := &handlers.SearchController{PageRepository: pageRepository}
	r.HandleFunc("/search", sc.ShowSearchResults).Methods("GET")
	r.HandleFunc("/api/search", sc.SearchAPI).Methods("GET") //returns json

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
