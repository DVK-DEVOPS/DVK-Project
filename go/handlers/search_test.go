// search_test.go
package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"DVK-Project/db"
)

// Mock repository
type MockPageRepo struct{}

func (m *MockPageRepo) FindSearchResults(searchStr string, language string) ([]db.Result, error) {
	if searchStr == "error" {
		return nil, errors.New("mock DB error")
	}
	return []db.Result{
		{Title: "Test Page", Url: "/test", Content: "Some content", Language: language},
	}, nil
}

func TestShowSearchResults(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		language       string
		expectContains string
		statusCode     int
	}{
		{
			name:           "Empty query triggers no results",
			query:          "",
			language:       "",
			expectContains: "No results",
			statusCode:     http.StatusOK,
		},
		{
			name:           "Normal query returns result",
			query:          "test",
			language:       "en",
			expectContains: "Test Page",
			statusCode:     http.StatusOK,
		},
		{
			name:           "Database error returns 500",
			query:          "error",
			language:       "",
			expectContains: "Database error",
			statusCode:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &SearchController{
				PageRepository: &MockPageRepo{},
				RenderTemplate: func(w http.ResponseWriter, r *http.Request, filename string, data interface{}) {
					if data == nil {
						w.Write([]byte("No results"))
					} else {
						w.Write([]byte("Test Page"))
					}
				},
			}

			req := httptest.NewRequest("GET", "/search?q="+tt.query+"&language="+tt.language, nil)
			w := httptest.NewRecorder()

			sc.ShowSearchResults(w, req)

			resp := w.Result()
			body := w.Body.String()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, resp.StatusCode)
			}
			if !strings.Contains(body, tt.expectContains) {
				t.Errorf("expected body to contain %q, got %q", tt.expectContains, body)
			}
		})
	}
}
