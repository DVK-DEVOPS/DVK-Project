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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of currently active users",
		},
	)
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		if path == "" {
			path = r.URL.Path
		}

		httpRequestsTotal.WithLabelValues(
			r.Method,
			path,
			strconv.Itoa(rw.statusCode),
		).Inc()

		httpRequestDuration.WithLabelValues(
			r.Method,
			path,
		).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {

	dsn := config.GetSentryDSN()
	if dsn != "" {
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

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

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
	r.HandleFunc("/password-reset", lh.ShowPasswordReset).Methods("GET")
	r.HandleFunc("/api/password-reset", lh.ResetPassword).Methods("POST")

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

	r.Use(metricsMiddleware)

	fmt.Println("Registered routes:") //debug
	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, _ := route.GetPathTemplate()
		fmt.Println(t)
		return nil
	}); err != nil {
		fmt.Printf("error walking routes: %v\n", err)
	}

	log.Println("Server running on :8080")
	log.Println("Metrics available at :8080/metrics")

	if dsn == "" {
		log.Fatal(http.ListenAndServe(":8080", r))
		return
	}
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         1 * time.Second,
	})
	log.Fatal(http.ListenAndServe(":8080", sentryHandler.Handle(r)))
}
