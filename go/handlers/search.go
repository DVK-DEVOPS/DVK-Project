package handlers

import (
	"DVK-Project/db"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type SearchController struct {
	PageRepository *db.PageRepository
}

// ### Page Templating ###
func (sc *SearchController) ShowSearchResults(w http.ResponseWriter, r *http.Request) {
	searchStr := r.FormValue("q")
	language := r.FormValue("language")

	if searchStr == "" { //trigger 'No results' block in html.
		tmpl := template.Must(template.ParseFiles("templates/search.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			fmt.Printf("search.go: failed to execute template: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	results, err := sc.PageRepository.FindSearchResults(searchStr, language)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	if err := tmpl.Execute(w, results); err != nil {
		fmt.Printf("search.go: failed to execute template: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// SearchAPI godoc
// @Summary      Search pages
// @Description  Search pages by query string (q) and optional language code
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        q         query     string  true   "Search query"
// @Param        language  query     string  false  "Language code (optional)"
// @Success      200 {object} db.Result"
// @Failure      422 {object} RequestValidationError "Validation Error - missing/invalid parameters"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /api/search [get]
func (sc *SearchController) SearchAPI(w http.ResponseWriter, req *http.Request) {
	searchStr := req.URL.Query().Get("q")
	language := req.URL.Query().Get("language")

	if searchStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		msg := "query parameter is required"
		validationErr := RequestValidationError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    &msg,
			Detail:     nil,
		}
		_ = json.NewEncoder(w).Encode(validationErr)
		return
	}

	results, err := sc.PageRepository.FindSearchResults(searchStr, language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": results,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type RequestValidationError struct {
	StatusCode int         `json:"statusCode"`
	Message    *string     `json:"message,omitempty"`
	Detail     interface{} `json:"detail,omitempty"` // can be string or null
}
