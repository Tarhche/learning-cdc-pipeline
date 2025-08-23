package searchindex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// TypesenseRepository implements the SearchIndexRepository interface
type TypesenseRepository struct {
	client *typesense.Client
	config TypesenseConfig
}

// TypesenseConfig holds the configuration for Typesense connection
type TypesenseConfig struct {
	Host   string
	Port   int
	APIKey string
}

// NewTypesenseRepository creates a new Typesense repository instance
func NewTypesenseRepository(config TypesenseConfig) *TypesenseRepository {
	return &TypesenseRepository{
		config: config,
	}
}

// Connect establishes a connection to Typesense
func (r *TypesenseRepository) Connect() error {
	serverURL := fmt.Sprintf("http://%s:%d", r.config.Host, r.config.Port)

	client := typesense.NewClient(
		typesense.WithServer(serverURL),
		typesense.WithAPIKey(r.config.APIKey),
	)

	// Test the connection by trying to retrieve collections
	ctx := context.Background()
	_, err := client.Collections().Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Typesense: %w", err)
	}

	r.client = client
	log.Printf("Connected to Typesense at %s:%d", r.config.Host, r.config.Port)
	return nil
}

// Close closes the Typesense connection
func (r *TypesenseRepository) Close() error {
	// Typesense client doesn't have a close method
	return nil
}

// CreateCollection creates a new collection in Typesense
func (r *TypesenseRepository) CreateCollection(schema map[string]interface{}) error {
	ctx := context.Background()

	// Check if collection already exists
	collections, err := r.client.Collections().Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve collections: %w", err)
	}

	collectionName := schema["name"].(string)
	for _, collection := range collections {
		if collection.Name == collectionName {
			log.Printf("Collection %s already exists", collectionName)
			return nil
		}
	}

	// Convert the schema map to CollectionSchema
	collectionSchema := &api.CollectionSchema{
		Name: collectionName,
	}

	// Add fields
	if fields, ok := schema["fields"].([]map[string]interface{}); ok {
		for _, field := range fields {
			fieldSchema := api.Field{
				Name: field["name"].(string),
				Type: field["type"].(string),
			}
			if optional, ok := field["optional"].(bool); ok && optional {
				fieldSchema.Optional = &optional
			}
			if facet, ok := field["facet"].(bool); ok {
				fieldSchema.Facet = &facet
			}
			if index, ok := field["index"].(bool); ok {
				fieldSchema.Index = &index
			}
			collectionSchema.Fields = append(collectionSchema.Fields, fieldSchema)
		}
	}

	// Add default sorting field
	if defaultSortingField, ok := schema["default_sorting_field"].(string); ok {
		collectionSchema.DefaultSortingField = &defaultSortingField
	}

	// Create the collection
	_, err = r.client.Collections().Create(ctx, collectionSchema)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	log.Printf("Created collection: %s", collectionName)
	return nil
}

// UpsertDocument upserts a document to a collection
func (r *TypesenseRepository) UpsertDocument(collectionName string, document interface{}) error {
	ctx := context.Background()

	// Upsert the document
	_, err := r.client.Collection(collectionName).Documents().Upsert(ctx, document)
	if err != nil {
		return fmt.Errorf("failed to upsert document: %w", err)
	}

	return nil
}

// DeleteDocument deletes a document from a collection
func (r *TypesenseRepository) DeleteDocument(collectionName string, documentID string) error {
	ctx := context.Background()

	// Delete the document
	_, err := r.client.Collection(collectionName).Document(documentID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

// SearchDocuments searches for documents in a collection
func (r *TypesenseRepository) SearchDocuments(collectionName, query string, searchParams map[string]interface{}) ([]interface{}, error) {
	ctx := context.Background()

	// Prepare search parameters
	searchParameters := &api.SearchCollectionParams{
		Q: query,
	}

	// Add additional search parameters if provided
	if filterBy, ok := searchParams["filter_by"].(string); ok {
		searchParameters.FilterBy = &filterBy
	}

	if sortBy, ok := searchParams["sort_by"].(string); ok {
		searchParameters.SortBy = &sortBy
	}

	if queryBy, ok := searchParams["query_by"].(string); ok {
		searchParameters.QueryBy = queryBy
	}

	// Add pagination parameters
	if page, ok := searchParams["page"].(int); ok {
		searchParameters.Page = &page
	}

	if perPage, ok := searchParams["per_page"].(int); ok {
		searchParameters.PerPage = &perPage
	}

	// Perform the search
	searchResult, err := r.client.Collection(collectionName).Documents().Search(ctx, searchParameters)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	// Convert results to interface slice
	var results []interface{}
	if searchResult.Hits != nil {
		for _, hit := range *searchResult.Hits {
			// Return the raw document with additional metadata
			document := *hit.Document
			document["_text_match"] = hit.TextMatch

			results = append(results, document)
		}
	}

	return results, nil
}

// GetAllDocuments retrieves all documents from a collection
func (r *TypesenseRepository) GetAllDocuments(collectionName string) ([]interface{}, error) {
	ctx := context.Background()

	// Export all documents from the collection
	exportReader, err := r.client.Collection(collectionName).Documents().Export(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to export documents: %w", err)
	}
	defer exportReader.Close()

	// Read the exported data
	exportData, err := io.ReadAll(exportReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read exported documents: %w", err)
	}

	// Parse the JSONL format (each line is a JSON document)
	lines := strings.Split(string(exportData), "\n")
	var results []interface{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var document map[string]interface{}
		if err := json.Unmarshal([]byte(line), &document); err != nil {
			// Skip invalid lines
			continue
		}
		results = append(results, document)
	}

	return results, nil
}
