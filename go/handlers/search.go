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
	searchStr := r.FormValue("query")

	results, err := sc.PageRepository.FindSearchResults(searchStr)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	tmpl.Execute(w, results)
}

func (sc *SearchController) SearchAPI(w http.ResponseWriter, req *http.Request) {
	searchStr := req.URL.Query().Get("q")

	results, err := sc.PageRepository.FindSearchResults(searchStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
