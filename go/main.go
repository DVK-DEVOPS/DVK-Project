package main

import (
	"DVK-Project/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/login", handlers.ShowLogin)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
