// The version package simply prints the version of the dashboard binary file.
package dashboard

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/browser"
	"github.com/senzing/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// DashboardImpl is the default implementation of the Dashboard interface.
type DashboardImpl struct {
	ServerPort int
	TtyOnly    bool
	logger     logging.LoggingInterface
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var content embed.FS

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func fileName(Uri string) string {
	return "static" + Uri
}

// ----------------------------------------------------------------------------
// Helper methods
// ----------------------------------------------------------------------------

func (dashboard *DashboardImpl) contentReadFileError(request *http.Request, err error) {
	dashboard.logger.Log(3001, request.RequestURI, request, err)
}

func (dashboard *DashboardImpl) writeResponse(responseWriter http.ResponseWriter, request *http.Request) {
	fileBytes, err := content.ReadFile(fileName(request.RequestURI))
	if err != nil {
		dashboard.contentReadFileError(request, err)
	}
	responseWriter.Write(fileBytes)
}

// ----------------------------------------------------------------------------
// Routers
// ----------------------------------------------------------------------------

func (dashboard *DashboardImpl) cssRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/css")
		dashboard.writeResponse(responseWriter, request)
	})
	return router
}

func (dashboard *DashboardImpl) pngRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/png")
		dashboard.writeResponse(responseWriter, request)
	})
	return router
}

func (dashboard *DashboardImpl) includeRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		dashboard.writeResponse(responseWriter, request)
	})
	return router
}

func (dashboard *DashboardImpl) jsRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/javascript")
		dashboard.writeResponse(responseWriter, request)
	})
	return router
}

func (dashboard *DashboardImpl) svgRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "image/svg+xml")
		dashboard.writeResponse(responseWriter, request)
	})
	return router
}

func (dashboard *DashboardImpl) webRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		dashboard.writeResponse(responseWriter, request)
	})
	return router
}

// ----------------------------------------------------------------------------
// Service
// ----------------------------------------------------------------------------

// Print a version string.
func (dashboard *DashboardImpl) Serve(ctx context.Context) error {
	dashboard.logger, _ = logging.New()

	// Print information for user.

	fmt.Printf("\n"+IdMessages[2003]+"\n\n", dashboard.ServerPort)

	// Log entry parameters.

	dashboard.logger.Log(2000, dashboard)

	// Set up a router to route http request.

	router := chi.NewRouter()

	// A good base middleware stack.

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	// Routes

	router.Mount("/css", dashboard.cssRouter())
	router.Mount("/js", dashboard.jsRouter())
	router.Mount("/png", dashboard.pngRouter())
	router.Mount("/include", dashboard.includeRouter())
	router.Mount("/svg", dashboard.svgRouter())
	router.Mount("/web", dashboard.webRouter())

	// Specific URIs.

	router.Get("/favicon.ico", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "image/x-icon")
		fileBytes, err := content.ReadFile(fileName("/img/favicon.ico"))
		if err != nil {
			dashboard.contentReadFileError(request, err)
		}
		responseWriter.Write(fileBytes)
	})

	router.Get("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		fileBytes, err := content.ReadFile(fileName("/web/index.html"))
		if err != nil {
			dashboard.contentReadFileError(request, err)
		}
		responseWriter.Write(fileBytes)
	})

	// Start a web browser.  Unless disabled.

	if !dashboard.TtyOnly {
		_ = browser.OpenURL(fmt.Sprintf("http://localhost:%d", dashboard.ServerPort))
	}

	// Start router.

	return http.ListenAndServe(":"+strconv.Itoa(dashboard.ServerPort), router)

}
