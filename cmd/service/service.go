// The version package simply prints the version of the go-fileindex binary file.
package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/docktermj/dashboard/box"
	"github.com/docktermj/dashboard/models/fileindex"
	log "github.com/docktermj/go-logger/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func viperAsJson() string {
	viperConfig := viper.AllSettings()
	viperByteArray, err := json.Marshal(viperConfig)
	if err != nil {
		log.Fatalf("Unable to marshal viper config to JSON: %v", err)
	}
	return string(viperByteArray)
}

// ----------------------------------------------------------------------------
// Routers
// ----------------------------------------------------------------------------

func apiRouter(database *sql.DB) chi.Router {

	handlerDatabase := &fileindex.DB{database}
	handler := &ApiHandler{
		datastore: handlerDatabase,
	}

	router := chi.NewRouter()
	router.Get("/", handler.root)
	router.Get("/count", handler.count)
	router.Get("/duplicates/sha256", handler.duplicatesSha256)
	router.Get("/everything", handler.everything)
	router.Get("/id/{id}", handler.id)
	router.Get("/modified/{modified}", handler.modified)
	router.Get("/name/{name}", handler.name)
	router.Get("/path/{path}", handler.path)
	router.Get("/sha256/{sha256}", handler.sha256)
	router.Get("/size/{size}", handler.size)
	router.Get("/unique/sha256", handler.uniqueSha256)
	router.Get("/volume/{volume}", handler.volume)
	return router
}

func cssRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/css")
		responseWriter.Write(box.Get(request.RequestURI))
	})
	return router
}

func imagesRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/png")
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
func Service(port string, sqliteFileName string) {
	// Perform command.

	// Prepare SQLite database.

	database, err := sql.Open("sqlite3", sqliteFileName)
	if err != nil {
		log.Errorf("%s cannot be opened. Error:  %v", sqliteFileName, err)
		return
	}

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

	router.Mount("/api", apiRouter(database))
	router.Mount("/css", cssRouter())
	router.Mount("/js", jsRouter())
	router.Mount("/images", imagesRouter())
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

	http.ListenAndServe(":"+port, router)
}

// ----------------------------------------------------------------------------
// Command pattern "Execute" function.
// ----------------------------------------------------------------------------

// The Command sofware design pattern's Execute() method.
func Execute() {

	// Print context parameters.
	// TODO: Formalize entry parameters

	log.Info(viperAsJson())

	// Get parameters from viper.

	var port = viper.GetString("port")
	var sqliteFileName = viper.GetString("sqlite_file_name")

	// Check parameters.

	errors := 0

	if port == "" {
		errors += 1
		log.Warn("--port not set.")
	}
	if sqliteFileName == "" {
		errors += 1
		log.Warn("--sqlite-file-name not set.")
	}

	// If any errors, exit.

	if errors > 0 {
		return
	}

	// Perform command.

	Service(port, sqliteFileName)

}
