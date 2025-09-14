package models

type AuthResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
}

type HTTPValidationError struct {
	Message string `json:"message"`
}
