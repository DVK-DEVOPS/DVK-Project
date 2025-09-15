package handlers

import (
	"DVK-Project/db"
	"encoding/json"
	"html/template"
	"net/http"
)

type SearchController struct {
	PageRepository *db.PageRepository
}

func (sc *SearchController) ShowSearchResults(w http.ResponseWriter, r *http.Request) {
	searchStr := r.FormValue("q")
	language := "en" //TODO: Add dropdown menu or selection to searchform

	results, err := sc.PageRepository.FindSearchResults(searchStr, language)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	tmpl.Execute(w, results)
}

// TODO: Change the Swagger annotations to reflect the OpenAPI spec. Code has already been changed to match.
// SearchAPI godoc
// @Summary      Search pages
// @Description  Search pages by query string
// @Tags         search
// @Accept       json
// @Produce      json
// @Param        query   query     string  true  "Search query"
// @Success 200 {object} db.Result
// @Failure      400     {object}  map[string]string "Bad Request"
// @Failure      500     {object}  map[string]string "Internal Server Error"
// @Router       /search [get]
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
