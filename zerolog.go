package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the logger
	log.Info().Msg("Starting the application")

	// Simulate an error
	log.Error().Str("error_type", "critical").Msg("An error occurred")

	// Log some structured data
	log.Info().Str("user", "john_doe").Int("age", 30).Msg("User login")

	// Shutdown the application
	log.Info().Msg("Shutting down the application")
}
