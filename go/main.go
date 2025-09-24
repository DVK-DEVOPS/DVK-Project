package main

import (
	"DVK-Project/client"
	"DVK-Project/config"
	"DVK-Project/db"
	"DVK-Project/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//debugging
	godotenv.Load()

	fmt.Println("API_KEY:", config.GetAPIKey())

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	userRepository := db.NewUserRepository(database)
	pageRepository := db.NewPageRepository(database)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/about", handlers.AboutHandler).Methods("GET")

	lh := &handlers.LoginHandler{UserRepository: userRepository}
	r.HandleFunc("/login", lh.ShowLogin).Methods("GET")
	r.HandleFunc("/api/login", lh.Login).Methods("POST")

	r.HandleFunc("/logout", lh.Logout).Methods("GET")

	rc := &handlers.RegistrationController{UserRepository: userRepository}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")

	sc := &handlers.SearchController{PageRepository: pageRepository}
	r.HandleFunc("/search", sc.ShowSearchResults).Methods("GET")
	r.HandleFunc("/api/search", sc.SearchAPI).Methods("GET") //returns json

	apiClient := client.NewAPIClient()
	wc := handlers.NewWeatherController(apiClient)
	r.HandleFunc("/weather", wc.ShowWeatherPage).Methods("GET")
	r.HandleFunc("/api/weather", wc.GetWeatherForecast).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
