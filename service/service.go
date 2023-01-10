// The version package simply prints the version of the dashboard binary file.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/docktermj/dashboard/box"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/spf13/viper"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ServiceImpl is the default implementation of the Service interface.
type HttpServerImpl struct {
	Port     int
	LogLevel logger.Level
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func viperAsJson() string {
	viperConfig := viper.AllSettings()
	viperByteArray, err := json.Marshal(viperConfig)
	if err != nil {
		fmt.Printf("Unable to marshal viper config to JSON: %v", err)
	}
	return string(viperByteArray)
}

// ----------------------------------------------------------------------------
// Routers
// ----------------------------------------------------------------------------

func cssRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/css")
		responseWriter.Write(box.Get(request.RequestURI))
	})
	return router
}

func pngRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/png")
		responseWriter.Write(box.Get(request.RequestURI))
	})
	return router
}

func includeRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write(box.Get(request.RequestURI))
	})
	return router
}

func jsRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/javascript")
		responseWriter.Write(box.Get(request.RequestURI))
	})
	return router
}

func svgRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "image/svg+xml")
		responseWriter.Write(box.Get(request.RequestURI))
	})
	return router
}

func webRouter() chi.Router {
	handler := &WebHandler{}
	router := chi.NewRouter()
	router.Get("/all", handler.all)
	router.Get("/duplicates/sha256", handler.duplicatesSha256)
	router.Get("/id/{id}", handler.id)
	router.Get("/modified/{modified}", handler.modified)
	router.Get("/name/{name}", handler.name)
	router.Get("/path/{path}", handler.path)
	router.Get("/scan", handler.scan)
	router.Get("/sha256/{sha256}", handler.sha256)
	router.Get("/size/{size}", handler.size)
	router.Get("/unique/sha256", handler.uniqueSha256)
	router.Get("/volume/{volume}", handler.volume)
	router.Get("/*", handler.staticHtml)
	return router
}

// ----------------------------------------------------------------------------
// Service
// ----------------------------------------------------------------------------

// Print a version string.
func (httpServer *HttpServerImpl) Serve(ctx context.Context) error {
	var err error = nil
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, httpServer.LogLevel)

	// Print information for user.

	fmt.Printf("\n"+IdMessages[2003]+"\n\n", httpServer.Port)

	// Log entry parameters.

	logger.Log(2000, httpServer)

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

	router.Mount("/css", cssRouter())
	router.Mount("/js", jsRouter())
	router.Mount("/png", pngRouter())
	router.Mount("/include", includeRouter())
	router.Mount("/svg", svgRouter())
	router.Mount("/web", webRouter())

	// Specific URIs.

	router.Get("/favicon.ico", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "image/x-icon")
		responseWriter.Write(box.Get("/img/favicon.ico"))
	})

	router.Get("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Write(box.Get("/web/index.html"))
	})

	// Start router.

	http.ListenAndServe(":"+strconv.Itoa(httpServer.Port), router)

	return err
}
