package client

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type APIClient struct {
	Client *http.Client
	Url    string
	ApiKey string
}

func NewAPIClient() *APIClient {
	_ = godotenv.Load()
	url := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")
	if url == "" || apiKey == "" {
		log.Fatal("API_URL or API_KEY not set")
	}
	return &APIClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Url:    url,
		ApiKey: apiKey,
	}
}
