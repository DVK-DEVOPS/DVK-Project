package handlers

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/getsentry/sentry-go"
)

var templatesFS embed.FS

func captureAndRespond(w http.ResponseWriter, r *http.Request, err error, msg string, code int) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		hub.CaptureException(err)
	}
	http.Error(w, msg, code)
}

func renderTemplate(w http.ResponseWriter, r *http.Request, filename string, data interface{}) {
	tmpl, err := template.ParseFS(templatesFS, "templates/"+filename)
	if err != nil {
		captureAndRespond(w, r, err, "Template not found", http.StatusNotFound)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		captureAndRespond(w, r, err, "Template rendering failed", http.StatusInternalServerError)
	}
}
