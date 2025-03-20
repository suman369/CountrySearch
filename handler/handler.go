package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"countrysearch/models"
	"countrysearch/service"
	"countrysearch/utils"
)

// Handler handles HTTP requests
type Handler struct {
	service *service.CountryService
	logger  utils.Logger
}

// NewHandler creates a new HTTP handler
func NewHandler(service *service.CountryService, logger utils.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// SetupRouter sets up the Gin router with all endpoints and middleware
func SetupRouter(service *service.CountryService, logger utils.Logger) *gin.Engine {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Use middleware
	router.Use(gin.Recovery())
	router.Use(LoggingMiddleware(logger))

	// Create handler
	handler := NewHandler(service, logger)

	// Define API routes
	api := router.Group("/api")
	{
		countries := api.Group("/countries")
		{
			countries.GET("/search", handler.SearchCountry)
		}
	}

	return router
}

// SearchCountry handles country search requests
func (h *Handler) SearchCountry(c *gin.Context) {
	// Get name parameter
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, models.Error{
			Message: "Name parameter is required",
			Status:  http.StatusBadRequest,
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Call service
	country, err := h.service.SearchCountry(ctx, name)
	if err != nil {
		h.logger.Error("Error searching for country", "error", err)
		c.JSON(http.StatusNotFound, models.Error{
			Message: "Country not found",
			Status:  http.StatusNotFound,
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, country)
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		logger.Info(
			"HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", latency,
			"user-agent", c.Request.UserAgent(),
		)
	}
}
