package handlers

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/getsentry/sentry-go"
)

//go:embed templates/*
var templatesFS embed.FS

func captureAndRespond(w http.ResponseWriter, r *http.Request, err error, msg string, code int) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		hub.CaptureException(err)
	}
	http.Error(w, msg, code)
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, filename string, data interface{}) {
	// ensure data is always a map for template injection
	renderData, ok := data.(map[string]interface{})
	if !ok || data == nil {
		renderData = map[string]interface{}{}
	}

	// auto-inject session info
	user, logged := GetUserFromSession(r)
	renderData["LoggedIn"] = logged
	if logged {
		renderData["User"] = user
	}

	tmpl, err := template.ParseFS(templatesFS, "templates/"+filename)
	if err != nil {
		captureAndRespond(w, r, err, "Template not found", http.StatusNotFound)
		return
	}

	if err := tmpl.Execute(w, renderData); err != nil {
		captureAndRespond(w, r, err, "Template rendering failed", http.StatusInternalServerError)
	}
}
