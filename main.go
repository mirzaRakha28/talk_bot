package main

import (
	"log"
	"net/http"
	"seatalk-bot/internal/config"
	"seatalk-bot/pkg/eventcallback"
)

func main() {
	// Load the configuration
	cfg := config.LoadConfig() // Ensure you have a LoadConfig function in your config package

	// Initialize the EventCallbackService
	eventService := eventcallback.NewEventCallbackService(cfg)

	// Set up the HTTP handler for event callbacks
	http.HandleFunc("/event-callback", eventService.HandleEventCallback)

	// Start the HTTP server
	port := ":8080" // Change the port as needed
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
