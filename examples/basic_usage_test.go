package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kylejryan/mocument/mock"
	"github.com/stretchr/testify/assert"
)

var mockDBClient *mock.MockDocDB
var mockSecretsClient *mock.MockSecretsManager

func init() {
	mockConfig := &mock.MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDBClient = mock.NewMockDocDB(mockConfig)

	mockSecretsClient = mock.NewMockSecretsManager()
	mockSecretsClient.AddSecret("MONGODB_SECRET_ID", "mongodb://localhost:27017")

	// Set the environment variable for the test
	os.Setenv("MONGODB_SECRET_ID", "MONGODB_SECRET_ID")
}

func mockHandler(ctx context.Context, event MyEvent) (string, error) {
	// Use the mock clients instead of the real AWS clients
	collection := mockDBClient

	// Insert a document
	doc := map[string]interface{}{"name": event.Name}
	err := collection.InsertDocument("collection", doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %w", err)
	}

	// Find the document
	filter := map[string]interface{}{"name": event.Name}
	results, err := collection.FindDocument("collection", filter)
	if err != nil {
		return "", fmt.Errorf("failed to find document: %w", err)
	}

	return fmt.Sprintf("Found documents: %+v", results), nil
}

func TestHandler(t *testing.T) {
	// Create a fake event
	event := MyEvent{Name: "test"}

	// Call the handler with the mock client
	result, err := mockHandler(context.Background(), event)
	assert.NoError(t, err)
	assert.Contains(t, result, "Found documents")
}
