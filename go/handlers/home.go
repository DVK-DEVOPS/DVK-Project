package handlers

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index.html", nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "about.html", nil)
}
