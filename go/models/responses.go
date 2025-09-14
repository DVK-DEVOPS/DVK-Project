package models

type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type HTTPValidationError struct {
	Detail []ValidationErrorDetail `json:"detail"`
}

type ValidationErrorDetail struct {
	Loc  []interface{} `json:"loc"`
	Msg  string        `json:"msg"`
	Type string        `json:"type"`
}
