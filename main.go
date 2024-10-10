package main

import (
	"log"
	"net/http"
)

func main() {
	// Load the configuration
	// cfg := config.LoadConfig() // Ensure you have a LoadConfig function in your config package

	// Initialize the EventCallbackService
	// eventService := eventcallback.NewEventCallbackService(cfg)

	// Set up the HTTP handler for event callbacks
	// http.HandleFunc("/event-callback", eventService.HandleEventCallback)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Start the HTTP server
	port := ":6969" // Change the port as needed
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
