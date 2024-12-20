package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"WASA/WASAText/pkg/db"
	"WASA/WASAText/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize router
	router := mux.NewRouter()

	// Register user routes
	routes.RegisterUserRoutes(router)

	// Start server
	go func() {
		log.Println("Server running on port 8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	waitForShutdown()
}

func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}
