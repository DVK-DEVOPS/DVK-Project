package models

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type APIClient struct {
	client *http.Client
	url    string
	apiKey string
}

func NewAPIClient() *APIClient {
	_ = godotenv.Load()
	url := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")
	if url == "" || apiKey == "" {
		log.Fatal("API_URL or API_KEY not set")
	}
	return &APIClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url:    url,
		apiKey: apiKey,
	}
}
