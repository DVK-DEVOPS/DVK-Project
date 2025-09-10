package handlers

import (
	"DVK-Project/db"
	"html/template"
	"net/http"
)

type SearchController struct {
}

func (sc *SearchController) ShowSearchResults(w http.ResponseWriter, r *http.Request) {
	searchStr := r.FormValue("query")

	results, err := db.FindSearchResults(searchStr)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	tmpl.Execute(w, results)
}
