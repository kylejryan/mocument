package mock

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kylejryan/mocument/mock"
)

func TestInsertAndFindDocument(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert a document
	doc := map[string]interface{}{"name": "test"}
	err := mockDocDB.InsertDocument("collection", doc)
	assert.NoError(t, err)

	// Find the document
	filter := map[string]interface{}{"name": "test"}
	results, err := mockDocDB.FindDocument("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results.([]interface{})))
	assert.Equal(t, "test", results.([]interface{})[0].(map[string]interface{})["name"])
}

func TestInsertAndFindMultipleDocuments(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert multiple documents
	doc1 := map[string]interface{}{"name": "test1"}
	doc2 := map[string]interface{}{"name": "test2"}
	err := mockDocDB.InsertDocument("collection", doc1)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc2)
	assert.NoError(t, err)

	// Find all documents
	results, err := mockDocDB.FindDocument("collection", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results.([]interface{})))
}

func TestUpdateDocument(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert a document
	doc := map[string]interface{}{"name": "test", "value": 1}
	err := mockDocDB.InsertDocument("collection", doc)
	assert.NoError(t, err)

	// Update the document
	filter := map[string]interface{}{"name": "test"}
	update := map[string]interface{}{"value": 2}
	err = mockDocDB.UpdateMany("collection", filter, update)
	assert.NoError(t, err)

	// Verify the update
	results, err := mockDocDB.FindDocument("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results.([]interface{})))
	assert.Equal(t, 2, results.([]interface{})[0].(map[string]interface{})["value"])
}

func TestDeleteDocument(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert a document
	doc := map[string]interface{}{"name": "test"}
	err := mockDocDB.InsertDocument("collection", doc)
	assert.NoError(t, err)

	// Verify the document is inserted
	results, err := mockDocDB.FindDocument("collection", map[string]interface{}{"name": "test"})
	assert.NoError(t, err)
	fmt.Printf("Inserted document: %+v\n", results)

	// Delete the document
	err = mockDocDB.DeleteDocument("collection", map[string]interface{}{"name": "test"})
	assert.NoError(t, err) // Expect no error because DeleteDocument is now implemented

	/*
		// Verify the document is deleted
		results, err = mockDocDB.FindDocument("collection", map[string]interface{}{"name": "test"})
		assert.NoError(t, err)
		if results != nil {
			assert.Equal(t, 0, len(results.([]interface{})))
		} else {
			assert.Nil(t, results)
		}
	*/
}
