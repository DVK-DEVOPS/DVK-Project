package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"

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
	// Ensure data is always a map for template injection
	renderData, ok := data.(map[string]interface{})
	if !ok || data == nil {
		renderData = map[string]interface{}{}
	}

	// Inject session info
	user, logged := GetUserFromSession(r)
	renderData["LoggedIn"] = logged
	if logged {
		renderData["User"] = user
	}

	// Parse nav partial + the requested page
	tmpl, err := template.ParseFS(templatesFS, "templates/nav.html", "templates/"+filename)
	if err != nil {
		captureAndRespond(w, r, err, "Template parsing failed", http.StatusInternalServerError)
		return
	}

	// Detect root template name from filename
	rootName := filepath.Base(filename)
	rootName = rootName[:len(rootName)-len(filepath.Ext(rootName))] // e.g., "search.html" → "search"

	// Execute the root template
	if err := tmpl.ExecuteTemplate(w, rootName, renderData); err != nil {
		captureAndRespond(w, r, err, "Template rendering failed", http.StatusInternalServerError)
	}
}
