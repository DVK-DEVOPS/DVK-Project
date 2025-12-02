package handlers

import (
	"DVK-Project/db"
	"DVK-Project/logging"
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type SearchController struct {
	//PageRepository *db.PageRepository
	PageRepository db.PageRepoInterface
	//RenderTemplate func(w http.ResponseWriter, r *http.Request, filename string, data interface{})
}

// ### Page Templating ###
func (sc *SearchController) ShowSearchResults(w http.ResponseWriter, r *http.Request) {
	searchStr := r.FormValue("q")
	language := r.FormValue("language")

	if searchStr == "" { //trigger 'No results' block in html.
		RenderTemplate(w, r, "search.html", map[string]interface{}{
			"Results": nil,
		})
		return
	}

	// Trace DB search
	span := sentry.StartSpan(r.Context(), "db.query",
		sentry.WithDescription("FindSearchResults"))
	results, err := sc.PageRepository.FindSearchResults(searchStr, language)
	span.Finish()
	if err != nil {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logging.Log.Info().
		Str("event", "search_performed").
		Str("query", searchStr).
		Msg("")

	RenderTemplate(w, r, "search.html", map[string]interface{}{
		"Results": results,
	})
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

	span := sentry.StartSpan(req.Context(), "db.query",
		sentry.WithDescription("FindSearchResults"))
	results, err := sc.PageRepository.FindSearchResults(searchStr, language)
	span.Finish()
	if err != nil {
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": results,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		if hub := sentry.GetHubFromContext(req.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type RequestValidationError struct {
	StatusCode int         `json:"statusCode"`
	Message    *string     `json:"message,omitempty"`
	Detail     interface{} `json:"detail,omitempty"` // can be string or null
}
