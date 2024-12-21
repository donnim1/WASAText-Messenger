/*
Package api exposes the main API engine for the texting application. All HTTP APIs are handled here, and business
logic should either be implemented in this package or delegated to specialized sub-packages for complex logic.

To use this package, create a new instance with `New()` by passing a valid `Config`. The resulting Router will provide
the `Handler()` function to return an HTTP handler suitable for use with an `http.Server` or middleware.

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

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()
*/
package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/donnim1/WASAText/service/database"
)

// Config provides dependencies and configuration for the API.
type Config struct {
	// Logger to record log entries
	Logger logrus.FieldLogger

	// Database instance to persist application data
	Database database.AppDatabase
}

// Router represents the API interface for handling HTTP requests.
type Router interface {
	// Handler returns the HTTP handler for this package's APIs
	Handler() http.Handler

	// Close releases any resources used by the package
	Close() error
}

// New creates and returns a new Router instance.
func New(cfg Config) (Router, error) {
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	// Create a new router for HTTP endpoints
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

	// baseLogger logs messages outside of request contexts (e.g., background tasks).
	baseLogger logrus.FieldLogger

	db database.AppDatabase
}

// Handler provides the HTTP handler for API endpoints.
func (r *_router) Handler() http.Handler {
	return r.router
}

// Close releases any resources used by the router.
func (r *_router) Close() error {
	// Close database connections or cleanup tasks here, if needed.
	return nil
}
