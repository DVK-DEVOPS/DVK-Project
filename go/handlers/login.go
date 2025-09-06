package handlers

import _ "github.com/gorilla/mux"

import "net/http"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Home Page",
		"Body":  "Hello from Gorilla Mux with HTML templates!",
	}
	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
