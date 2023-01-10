// The version package simply prints the version of the dashboard binary file.
package dashboard

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var content embed.FS

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

func fileName(Uri string) string {
	return "static" + Uri
}

// ----------------------------------------------------------------------------
// Routers
// ----------------------------------------------------------------------------

func cssRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/css")
		fileBytes, err := content.ReadFile(fileName(request.RequestURI))
		if err != nil {
			fmt.Printf(">>>> Error reading file: %s\n", fileName(request.RequestURI))
		}
		responseWriter.Write(fileBytes)
	})
	return router
}

func pngRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/png")
		fileBytes, err := content.ReadFile(fileName(request.RequestURI))
		if err != nil {
			fmt.Printf(">>>> Error reading file: %s\n", fileName(request.RequestURI))
		}
		responseWriter.Write(fileBytes)
	})
	return router
}

func includeRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		fileBytes, err := content.ReadFile(fileName(request.RequestURI))
		if err != nil {
			fmt.Printf(">>>> Error reading file: %s\n", fileName(request.RequestURI))
		}
		responseWriter.Write(fileBytes)
	})
	return router
}

func jsRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/javascript")
		fileBytes, err := content.ReadFile(fileName(request.RequestURI))
		if err != nil {
			fmt.Printf(">>>> Error reading file: %s\n", fileName(request.RequestURI))
		}
		responseWriter.Write(fileBytes)
	})
	return router
}

func svgRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "image/svg+xml")
		fileBytes, err := content.ReadFile(fileName(request.RequestURI))
		if err != nil {
			fmt.Printf(">>>> Error reading file: %s\n", fileName(request.RequestURI))
		}
		responseWriter.Write(fileBytes)
	})
	return router
}

func webRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		fileBytes, err := content.ReadFile(fileName(request.RequestURI))
		if err != nil {
			fmt.Printf(">>>> Error reading file: %s\n", fileName(request.RequestURI))
		}
		responseWriter.Write(fileBytes)
	})
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
		fileBytes, err := content.ReadFile("static/img/favicon.ico")
		if err != nil {
			fmt.Printf(">>>>>>>>>>>> Error reading: %s\n", request.RequestURI)
		}
		responseWriter.Write(fileBytes)
	})

	router.Get("/", func(responseWriter http.ResponseWriter, request *http.Request) {

		fileBytes, err := content.ReadFile("static/web/index.html")
		if err != nil {
			fmt.Printf(">>>>>>>>>>>> Error reading: %s\n", request.RequestURI)
		}
		responseWriter.Write(fileBytes)
	})

	// Start router.

	http.ListenAndServe(":"+strconv.Itoa(httpServer.Port), router)

	return err
}
