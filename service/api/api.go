/*
Package api exposes the main API engine. All HTTP APIs are handled here – so-called "business logic" should be here, or
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

Note: This project is built with Go 1.17.
*/
package api

import (
	"errors"
	"net/http"

	"github.com/donnim1/WASAText/service/api/handlers" // assumed package for endpoint handlers
	"github.com/donnim1/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database.AppDatabase where data are saved
	Database database.AppDatabase
}

// Router is the package API interface representing an API handler builder.
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package.
	Handler() http.Handler

	// Close terminates any resource used in the package.
	Close() error
}

type routerImpl struct {
	router     *httprouter.Router
	baseLogger logrus.FieldLogger
	db         database.AppDatabase
}

// New returns a new Router instance.
func New(cfg Config) (Router, error) {
	// Check if the configuration is correct.
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	// Create a new router where we will register HTTP endpoints.
	r := httprouter.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// Register endpoints.
	// Example endpoints: login, update username, and get conversations.
	// The handler functions below are assumed to be defined in the handlers package.
	r.POST("/session", handlers.HandleLogin)
	r.PUT("/user/username", handlers.HandleSetMyUserName)
	r.GET("/conversations", handlers.HandleGetMyConversations)

	// (Additional endpoints can be registered following the same pattern:
	// r.GET("/conversations/:conversationId", handlers.HandleGetConversation)
	// r.POST("/conversations/:conversationId/messages", handlers.HandleSendMessage)
	// etc.)

	return &routerImpl{
		router:     r,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}

// Handler returns an HTTP handler for the registered APIs.
func (r *routerImpl) Handler() http.Handler {
	return r.router
}

// Close terminates any resources used by the API package.
// (For now, there might not be any additional resources to close.)
func (r *routerImpl) Close() error {
	// Add cleanup logic if needed.
	return nil
}
