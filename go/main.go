package main

import (
	"DVK-Project/client"
	"DVK-Project/db"
	"DVK-Project/handlers"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total number of HTTP requests"},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "Duration of HTTP requests in seconds", Buckets: prometheus.DefBuckets},
		[]string{"method", "endpoint"},
	)
	activeUsers = promauto.NewGauge(prometheus.GaugeOpts{Name: "active_users", Help: "Number of currently active users"})
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		go func() {
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			if path == "" {
				path = r.URL.Path
			}

			httpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rw.statusCode)).Inc()
			httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		}()
	})
}

func main() {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	userRepo := db.NewUserRepository(dbConn)
	pageRepo := db.NewPageRepository(dbConn)

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
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

	r.Use(metricsMiddleware)

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
