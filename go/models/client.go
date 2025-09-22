package models

import "net/http"

type APIClient struct {
	client *http.Client
	url    string
	apiKey string
}
