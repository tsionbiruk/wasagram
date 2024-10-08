/*
Package api exposes the main API engine. All HTTP APIs are handled here - so-called "business logic" should be here, or
in a dedicated package (if that logic is complex enough).

To use this package, you should create a new instance with New() passing a valid Config. The resulting Router will have
the Router.Handler() function that returns a handler that can be used in a http.Server (or in other middlewares).

Example:

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: appdb,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	// ... other stuff here, like middleware chaining, etc.

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()

See the `main.go` file inside the `cmd/webapi` for a full usage example.
*/
package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/tsionbiruk/wasagram/service/database"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database.AppDatabase where data are saved
	Database database.AppDatabase
}

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Close terminates any resource used in the package
	Close() error
}

// New returns a new Router instance
func New(cfg Config) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}

type _router struct {
	router *httprouter.Router

	// baseLogger is a logger for non-requests contexts, like goroutines or background tasks not started by a request.
	// Use context logger if available (e.g., in requests) instead of this logger.
	baseLogger logrus.FieldLogger

	db database.AppDatabase
}

func (rt *_router) Authorize(w http.ResponseWriter, r *http.Request, username string) bool {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := strconv.ParseInt(tokenString, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse photo ID: %v", err.Error()), http.StatusInternalServerError)
		return false
	}
	// Validate the token by querying into the db

	var user_token int64
	user_token, err = rt.db.Gettoken(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User doesnt exist in Tokens table please log back in to get a token: %v", err.Error()), http.StatusInternalServerError)
		return false
	}
	if token != user_token {
		http.Error(w, "invalid token", http.StatusInternalServerError)
		return false
	}

	var tokentime time.Time
	tokentime, err = rt.db.Gettokentime(username, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to get token time please log back in to get a token time: %v", err.Error()), http.StatusInternalServerError)
		return false
	}

	if rt.db.Istokenexpired(tokentime) {
		http.Error(w, "Token expired. Please log back in to get a new token.", http.StatusUnauthorized)
		return false
	}

	return true

}
