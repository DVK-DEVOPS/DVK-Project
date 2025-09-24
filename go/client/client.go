package client

import (
	"fmt"
	"io"
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
func (c *APIClient) FetchForecast(city string) ([]byte, error) {
	fmt.Println("Fetching forecast")
	reqURL := fmt.Sprintf("%s?q=%s&appid=%s", c.Url, city, c.ApiKey)
	fmt.Println(reqURL)
	resp, err := c.Client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("Fetch error: %w", err)
	}
	fmt.Println("response status:", resp.Status)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Fetch error: %w", err)
	}
	return body, nil
}
