package apihandler

import (
	_ "embed"
	"html/template"
	"net/http"
	"templify/pkg/server/handler"
)

// embed the openapi JSON and HTML file into the binary
// so we can serve them without reading from the filesystem

//go:embed embedded/openapi.json
var openapiJSON []byte

//go:embed embedded/stoplight.html
var openapiHTMLStoplight []byte

// Get status
// (GET /info/status)
func (ah *APIHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	handler.RespondWithJSON(w, r, http.StatusOK, map[string]string{"status": "HEALTHY"})
}

// Get version
// (GET /info/version)
func (ah *APIHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	handler.RespondWithJSON(w, r, http.StatusOK, ah.Info)
}

// Get openapi JSON
// (GET /info/openapi.json)
func (ah *APIHandler) GetOpenAPIJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(openapiJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get openapi HTML
// (GET /info/openapi.html)
func (ah *APIHandler) GetOpenAPIHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var templateString = string(openapiHTMLStoplight)

	t, err := template.New("openapi").Parse(templateString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// replace the base URL in the HTML file
	// with the actual base URL of the server
	// and render to the response writer
	err = t.Execute(w, map[string]string{
		"BaseURL": ah.BaseURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
