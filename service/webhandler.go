// Support for the /web URL route.
package service

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/docktermj/dashboard/box"
	log "github.com/docktermj/go-logger/logger"
	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

// ----------------------------------------------------------------------------
// WebHandler structures
// ----------------------------------------------------------------------------

type WebHandler struct {
}

type FileIndexTableComplete struct {
	ApiUrl string
	Title  string
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func (handler *WebHandler) renderFileIndexTableComplete(responseWriter http.ResponseWriter, data FileIndexTableComplete) {

	// HTML headers.

	responseWriter.Header().Set("Content-Type", "text/html")

	// Render template.

	templateString := box.Get("/template/fileindex-table-complete.html")
	templateParsed, err := template.New("mike").Parse(string(templateString))
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.sha256() Error:  %v", err)
		return
	}
	err = templateParsed.Execute(responseWriter, data)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.sha256() Error:  %v", err)
		return
	}
}

// ----------------------------------------------------------------------------
// HTTP handlers.
//   Signature:  method(net/http.ResponseWriter, *net/http.Request)
// ----------------------------------------------------------------------------

func (handler *WebHandler) all(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: "/api/everything",
		Title:  "Everything",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) duplicatesSha256(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: "/api/duplicates/sha256",
		Title:  "Duplicates of SHA256",
	}

	// HTML headers.

	responseWriter.Header().Set("Content-Type", "text/html")

	// Render template.

	templateString := box.Get("/template/duplicates-sha256.html")
	templateParsed, err := template.New("mike").Parse(string(templateString))
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.sha256() Error:  %v", err)
		return
	}
	err = templateParsed.Execute(responseWriter, data)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.sha256() Error:  %v", err)
		return
	}
}

func (handler *WebHandler) scan(responseWriter http.ResponseWriter, request *http.Request) {

	// Input to template.

	data := struct {
		ApiUrl  string
		Title   string
		Volumes []string
	}{
		ApiUrl:  "/api/scan-chooser",
		Title:   "Scan",
		Volumes: []string{"a1", "b2", "c3"},
	}

	// HTML headers.

	responseWriter.Header().Set("Content-Type", "text/html")

	// Render template.

	templateString := box.Get("/template/scan.html")
	templateParsed, err := template.New("mike").Parse(string(templateString))
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.sha256() Error:  %v", err)
		return
	}
	err = templateParsed.Execute(responseWriter, data)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		log.Errorf("handler.sha256() Error:  %v", err)
		return
	}
}

func (handler *WebHandler) uniqueSha256(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: "/api/unique/sha256",
		Title:  "Unique of SHA256",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

// ----------------------------------------------------------------------------
// HTTP handlers. So called "byXxx"
//   Signature:  method(net/http.ResponseWriter, *net/http.Request)
// ----------------------------------------------------------------------------

func (handler *WebHandler) id(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/id/%s", chi.URLParam(request, "id")),
		Title:  "ID",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) modified(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/modified/%s", chi.URLParam(request, "modified")),
		Title:  "Modified",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) name(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/name/%s", chi.URLParam(request, "name")),
		Title:  "Name",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) path(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/path/%s", chi.URLParam(request, "path")),
		Title:  "Path",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) sha256(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/sha256/%s", chi.URLParam(request, "sha256")),
		Title:  "SHA256",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) size(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/size/%s", chi.URLParam(request, "size")),
		Title:  "Size",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}

func (handler *WebHandler) staticHtml(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write(box.Get(request.RequestURI))
}

func (handler *WebHandler) volume(responseWriter http.ResponseWriter, request *http.Request) {
	data := FileIndexTableComplete{
		ApiUrl: fmt.Sprintf("/api/volume/%s", chi.URLParam(request, "volume")),
		Title:  "Volume",
	}
	handler.renderFileIndexTableComplete(responseWriter, data)
}
