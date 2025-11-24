package main

import (
	"DVK-Project/client"
	"DVK-Project/db"
	"DVK-Project/handlers"
	"DVK-Project/metrics"
	"DVK-Project/scraper"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	userRepo := db.NewUserRepository(dbConn)
	pageRepo := db.NewPageRepository(dbConn)
	repo := db.NewPageRepository(dbConn)

	go scraper.Run("https://developer.mozilla.org/en-US/docs/Web/CSS/flex", repo)
	go scraper.Run("https://go.dev/doc/effective_go", repo)

	r := mux.NewRouter()

	r.Handle("/metrics", metrics.Handler())
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/about", handlers.AboutHandler).Methods("GET")

	lh := &handlers.LoginHandler{UserRepository: userRepo}
	r.HandleFunc("/login", lh.ShowLogin).Methods("GET")
	r.HandleFunc("/api/login", lh.Login).Methods("POST")
	r.HandleFunc("/password-reset", lh.ShowPasswordReset).Methods("GET")
	r.HandleFunc("/api/password-reset", lh.ResetPassword).Methods("POST")

	r.HandleFunc("/logout", lh.Logout).Methods("GET")

	rc := &handlers.RegistrationController{UserRepository: userRepo}
	r.HandleFunc("/register", rc.ShowRegistrationPage).Methods("GET")
	r.HandleFunc("/api/register", rc.Register).Methods("POST")

	sc := &handlers.SearchController{PageRepository: pageRepo}
	r.HandleFunc("/search", sc.ShowSearchResults).Methods("GET")
	r.HandleFunc("/api/search", sc.SearchAPI).Methods("GET")

	apiClient := client.NewAPIClient()
	wc := handlers.NewWeatherController(apiClient)
	r.HandleFunc("/weather", wc.ShowWeatherPage).Methods("GET")
	r.HandleFunc("/api/weather", wc.GetWeatherForecast).Methods("GET")

	r.Use(metrics.Middleware)

	fmt.Println("Registered routes:")
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, _ := route.GetPathTemplate()
		fmt.Println(t)
		return nil
	})

	log.Println("Server running on :8080")
	log.Println("Metrics available at :8080/metrics")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
