package main

import (
	"DVK-Project/client"
	"DVK-Project/config"
	"DVK-Project/db"
	"DVK-Project/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
)

func main() {

	// Initialize Sentry using configuration/env (only if DSN is provided)
	dsn := config.GetSentryDSN()
	if dsn != "" {
		// parse trace sample rate with conservative default
		traceRate := 0.05
		if v := os.Getenv("SENTRY_TRACES_SAMPLE_RATE"); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				traceRate = f
			}
		}
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			Environment:      config.GetSentryEnvironment(),
			Release:          os.Getenv("SENTRY_RELEASE"),
			Debug:            os.Getenv("SENTRY_DEBUG") == "true",
			EnableLogs:       true,
			TracesSampleRate: traceRate,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if event == nil || event.Request == nil {
					return event
				}
				// Redact cookies and authorization headers
				if event.Request.Headers != nil {
					if _, ok := event.Request.Headers["Authorization"]; ok {
						event.Request.Headers["Authorization"] = "[REDACTED]"
					}
					if _, ok := event.Request.Headers["Cookie"]; ok {
						event.Request.Headers["Cookie"] = "[REDACTED]"
					}
				}
				if event.Request.Cookies != "" {
					event.Request.Cookies = "[REDACTED]"
				}
				return event
			},
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	userRepository := db.NewUserRepository(database)
	pageRepository := db.NewPageRepository(database)

	r := mux.NewRouter()

	//Healthcheck
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

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

	fmt.Println("Registered routes:") //debug
	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, _ := route.GetPathTemplate()
		fmt.Println(t)
		return nil
	}); err != nil {
		fmt.Printf("error walking routes: %v\n", err)
	}

	log.Println("Server running on :8080")
	if dsn == "" {
		log.Fatal(http.ListenAndServe(":8080", r))
		return
	}
	// Wrap router with Sentry HTTP handler to capture panics and request context
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         1 * time.Second,
	})
	log.Fatal(http.ListenAndServe(":8080", sentryHandler.Handle(r)))
}
