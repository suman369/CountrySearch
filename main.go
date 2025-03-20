// main.go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"countrysearch/handler"
	"countrysearch/cache"
	"countrysearch/client"
	"countrysearch/configs"
	"countrysearch/service"
	"countrysearch/utils"
)

func main() {
	// Initialize logger
	logger := utils.NewLogger()
	logger.Info("Starting Country Search API")

	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded", "port", cfg.Port, "timeout", cfg.Timeout)

	// Initialize cache
	countryCache := cache.NewInMemoryCache()
	logger.Info("Cache initialized")

	// Initialize HTTP client
	httpClient := client.NewCountryClient(cfg.RestCountriesBaseURL, time.Duration(cfg.Timeout)*time.Second)
	logger.Info("HTTP client initialized")

	// Initialize service
	countryService := service.NewCountryService(httpClient, countryCache, logger)
	logger.Info("Country service initialized")

	// Initialize API handlers
	router := handler.SetupRouter(countryService, logger)
	logger.Info("API router initialized")

	// Initialize server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}
