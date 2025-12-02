package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

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
	// Always build a map for template injection
	renderData := map[string]interface{}{}

	// If caller already passed a map, copy it in
	if dataMap, ok := data.(map[string]interface{}); ok && dataMap != nil {
		for k, v := range dataMap {
			renderData[k] = v
		}
	} else if data != nil {
		// Otherwise, wrap the data under a generic key
		renderData["Data"] = data
	}

	// Inject session info
	user, logged := GetUserFromSession(r)
	renderData["LoggedIn"] = logged
	if logged {
		renderData["User"] = user
	}
	renderData["CurrentYear"] = time.Now().Year()

	// Parse nav partial + the requested page
	tmpl, err := template.ParseFS(templatesFS, "templates/nav.html", "templates/footer.html", "templates/"+filename)
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
