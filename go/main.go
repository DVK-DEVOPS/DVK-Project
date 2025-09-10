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
	authRepository := db.NewAuthRepository(database)
	userRepository := db.NewUserRepository(database)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)

	lh := &handlers.LoginHandler{AuthRepo: authRepository}
	r.HandleFunc("/login", lh.ShowLogin).Methods("GET")
	r.HandleFunc("/api/login", lh.Login).Methods("POST")

	rc := &handlers.RegistrationController{UserRepository: userRepository}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
