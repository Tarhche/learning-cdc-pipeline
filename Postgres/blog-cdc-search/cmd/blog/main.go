package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"blog-cdc-search/application/service"
	"blog-cdc-search/domain"
	"blog-cdc-search/infrastructure/database"
	"blog-cdc-search/infrastructure/repository"
	"blog-cdc-search/infrastructure/searchindex"
	"blog-cdc-search/infrastructure/web"
)

// SearchServiceAdapter adapts the actual SearchService to the interface
type SearchServiceAdapter struct {
	service *service.SearchService
}

func (a *SearchServiceAdapter) SearchPosts(ctx context.Context, params interface{}) (interface{}, error) {
	// Convert interface{} to SearchParams
	searchParams, ok := params.(service.SearchParams)
	if !ok {
		// Try to convert from map[string]interface{} if that's what was passed
		if paramMap, ok := params.(map[string]interface{}); ok {
			searchParams = service.SearchParams{
				Query:    paramMap["query"].(string),
				Page:     paramMap["page"].(int),
				PerPage:  paramMap["per_page"].(int),
				SortBy:   paramMap["sort_by"].(string),
				FilterBy: paramMap["filter_by"].(string),
			}
		} else {
			// Return error for unsupported type
			return nil, fmt.Errorf("unsupported params type: %T", params)
		}
	}

	return a.service.SearchPosts(ctx, searchParams)
}

func (a *SearchServiceAdapter) GetAllPostsFromIndex(ctx context.Context) ([]*domain.Post, error) {
	return a.service.GetAllPostsFromIndex(ctx)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "bloguser")
	dbPassword := getEnv("DB_PASSWORD", "blogpass")
	dbName := getEnv("DB_NAME", "blog")
	port := getEnv("PORT", "8085")

	// Typesense configuration
	typesenseHost := getEnv("TYPESENSE_HOST", "localhost")
	typesensePortStr := getEnv("TYPESENSE_PORT", "8108")
	typesensePort := 8108 // default
	if port, err := strconv.Atoi(typesensePortStr); err == nil {
		typesensePort = port
	}
	typesenseAPIKey := getEnv("TYPESENSE_API_KEY", "xyz")

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Blog Application...")

	// Connect to database
	db, err := database.NewPostgreSQLConnection(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL database")

	// Initialize repositories
	postRepo := repository.NewPostgreSQLPostRepository(db)

	// Initialize Typesense repository
	typesenseConfig := searchindex.TypesenseConfig{
		Host:   typesenseHost,
		Port:   typesensePort,
		APIKey: typesenseAPIKey,
	}
	typesenseRepo := searchindex.NewTypesenseRepository(typesenseConfig)

	// Connect to Typesense
	if err := typesenseRepo.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to Typesense: %v", err)
		log.Println("Search functionality will not be available")
	} else {
		log.Printf("Connected to Typesense at %s:%d", typesenseHost, typesensePort)
	}

	// Initialize services
	postService := service.NewPostService(postRepo)
	searchService := service.NewSearchService(typesenseRepo)

	// Create search service adapter
	searchServiceAdapter := &SearchServiceAdapter{
		service: searchService,
	}

	// Initialize handlers
	handlers := web.NewHandlers(postService, searchServiceAdapter)
	log.Printf("Handlers initialized: %+v", handlers)

	// Setup routes
	router := web.SetupRoutes(handlers)
	log.Printf("Routes setup complete")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
