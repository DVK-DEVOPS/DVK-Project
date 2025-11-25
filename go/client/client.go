package client

import (
	"DVK-Project/config"
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
	apiKey := config.GetSecret("API_KEY", os.Getenv("AZURE_KEYVAULT_NAME"), os.Getenv("API_KEY_SECRET_NAME"))
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

	resp, err := c.Client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	fmt.Println("response status:", resp.Status)
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("Warning client.go: failed to close response body: %v\n", cerr)
		}
	}()
	//defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	return body, nil
}
