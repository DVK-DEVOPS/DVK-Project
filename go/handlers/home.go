package handlers

import (
	"html/template"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureException(err)
		}
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}
}
