package mock

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kylejryan/mocument/mock"
)

func TestInsertAndFindDocument(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert a document
	doc := Document{"name": "test"}
	err := mockDocDB.InsertDocument("collection", doc)
	assert.NoError(t, err)

	// Find the document
	filter := Document{"name": "test"}
	results, err := mockDocDB.FindDocument("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "test", results[0]["name"])
}

func TestInsertAndFindMultipleDocuments(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert multiple documents
	doc1 := Document{"name": "test1"}
	doc2 := Document{"name": "test2"}
	err := mockDocDB.InsertDocument("collection", doc1)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc2)
	assert.NoError(t, err)

	// Find all documents
	results, err := mockDocDB.FindDocument("collection", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))
}

func TestUpdateDocument(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert a document
	doc := Document{"name": "test", "value": 1}
	err := mockDocDB.InsertDocument("collection", doc)
	assert.NoError(t, err)

	// Update the document using $set operator
	filter := Document{"name": "test"}
	update := Document{"$set": Document{"value": 2}}
	err = mockDocDB.UpdateMany("collection", filter, update)
	assert.NoError(t, err)

	// Verify the update
	results, err := mockDocDB.FindDocument("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 2, results[0]["value"])
}

func TestDeleteDocument(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert a document and capture the inserted document with _id
	doc := Document{"name": "test"}
	err := mockDocDB.InsertDocument("collection", doc)
	assert.NoError(t, err)

	// Find the inserted document to get its _id
	results, err := mockDocDB.FindDocument("collection", Document{"name": "test"})
	assert.NoError(t, err)
	fmt.Printf("Inserted document: %+v\n", results)
	assert.Equal(t, 1, len(results))
	insertedDoc := results[0]

	// Delete the document using its _id
	err = mockDocDB.DeleteDocument("collection", Document{"_id": insertedDoc["_id"]})
	assert.NoError(t, err)

	// Verify the document is deleted
	results, err = mockDocDB.FindDocument("collection", Document{"name": "test"})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))
}

func loadJSONFixture(filePath string, t *testing.T) Document {
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	var doc Document
	err = json.Unmarshal(data, &doc)
	assert.NoError(t, err)

	return doc
}

func TestInsertAndFindTransaction(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Load transaction from JSON fixture
	doc := loadJSONFixture("testdata/sample_transaction.json", t)

	// Insert the transaction document
	err := mockDocDB.InsertDocument("transactions", doc)
	assert.NoError(t, err)

	// Find the transaction document
	filter := Document{"ID": "txn001"}
	results, err := mockDocDB.FindDocument("transactions", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "txn001", results[0]["ID"])
	assert.Equal(t, 150.75, results[0]["Amount"])
	assert.Equal(t, "USD", results[0]["Currency"])
	assert.Equal(t, "Pending", results[0]["Status"])

	// Verify the items in the transaction
	items := results[0]["Items"].([]interface{})
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "prod001", items[0].(map[string]interface{})["ProductID"])
	assert.Equal(t, 50.25, items[0].(map[string]interface{})["Price"])
	assert.Equal(t, 2.0, items[1].(map[string]interface{})["Quantity"])
}

func TestInsertManyAndFindDocuments(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert multiple documents
	err := mockDocDB.InsertMany("collection", []interface{}{
		Document{"name": "test1", "value": 1},
		Document{"name": "test2", "value": 2},
		Document{"name": "test3", "value": 3},
	})
	assert.NoError(t, err)

	// Find all documents
	results, err := mockDocDB.FindDocument("collection", nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(results))

	// Validate that all documents are correctly inserted
	expectedValues := []int{1, 2, 3}
	for i, result := range results {
		assert.Equal(t, expectedValues[i], int(result["value"].(float64)))
	}

	// Filter and find a specific document
	filter := Document{"name": "test2"}
	specificResult, err := mockDocDB.FindDocument("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(specificResult))
	assert.Equal(t, "test2", specificResult[0]["name"])
	assert.Equal(t, 2, int(specificResult[0]["value"].(float64)))
}

func TestCountDocuments(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert multiple documents
	doc1 := Document{"name": "test1", "value": 1}
	doc2 := Document{"name": "test2", "value": 2}
	doc3 := Document{"name": "test3", "value": 3}
	err := mockDocDB.InsertDocument("collection", doc1)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc2)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc3)
	assert.NoError(t, err)

	// Count all documents in the collection
	count, err := mockDocDB.CountDocuments("collection", nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	// Count documents with a filter
	filter := Document{"value": 2}
	filteredCount, err := mockDocDB.CountDocuments("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, filteredCount)

	// Count documents with a filter that matches no documents
	filter = Document{"value": 99}
	noMatchCount, err := mockDocDB.CountDocuments("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 0, noMatchCount)
}

func TestUpdateOne(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert multiple documents
	doc1 := Document{"name": "test1", "value": 1}
	doc2 := Document{"name": "test2", "value": 2}
	doc3 := Document{"name": "test3", "value": 3}
	err := mockDocDB.InsertDocument("collection", doc1)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc2)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc3)
	assert.NoError(t, err)

	// Update one document
	filter := Document{"name": "test2"}
	update := Document{"value": 22}
	err = mockDocDB.UpdateOne("collection", filter, update)
	assert.NoError(t, err)

	// Verify the update
	results, err := mockDocDB.FindDocument("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 22, int(results[0]["value"].(float64)))

	// Verify other documents are not updated
	otherFilter := Document{"name": "test1"}
	otherResults, err := mockDocDB.FindDocument("collection", otherFilter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(otherResults))
	assert.Equal(t, 1, int(otherResults[0]["value"].(float64)))
}

func TestDeleteMany(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert multiple documents
	doc1 := Document{"name": "test1", "value": 1}
	doc2 := Document{"name": "test2", "value": 2}
	doc3 := Document{"name": "test2", "value": 3}
	doc4 := Document{"name": "test3", "value": 4}
	err := mockDocDB.InsertDocument("collection", doc1)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc2)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc3)
	assert.NoError(t, err)
	err = mockDocDB.InsertDocument("collection", doc4)
	assert.NoError(t, err)

	// Delete documents with name "test2"
	filter := Document{"name": "test2"}
	deletedCount, err := mockDocDB.DeleteMany("collection", filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, deletedCount)

	// Verify the remaining documents
	results, err := mockDocDB.FindDocument("collection", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))

	// Verify the correct documents were deleted
	for _, result := range results {
		docName := result["name"]
		assert.NotEqual(t, "test2", docName)
	}
}

func TestFindDocumentWithComplexFilter(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert documents
	docs := []Document{
		{"name": "Alice", "age": 30, "city": "New York"},
		{"name": "Bob", "age": 25, "city": "Los Angeles"},
		{"name": "Charlie", "age": 35, "city": "New York"},
	}
	for _, doc := range docs {
		err := mockDocDB.InsertDocument("users", doc)
		assert.NoError(t, err)
	}

	// Test complex filter
	filter := Document{
		"age":  Document{"$gt": 28},
		"city": "New York",
	}
	results, err := mockDocDB.FindDocument("users", filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))

	// Verify results
	for _, doc := range results {
		age := doc["age"].(float64)
		assert.True(t, age > 28)
		assert.Equal(t, "New York", doc["city"])
	}
}

func TestUpdateManyWithComplexUpdate(t *testing.T) {
	mockConfig := &MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := NewMockDocDB(mockConfig)

	// Insert documents
	docs := []Document{
		{"name": "Product A", "price": 100, "stock": 50},
		{"name": "Product B", "price": 200, "stock": 30},
		{"name": "Product C", "price": 150, "stock": 0},
	}
	for _, doc := range docs {
		err := mockDocDB.InsertDocument("products", doc)
		assert.NoError(t, err)
	}

	// Update documents with complex update
	filter := Document{"stock": Document{"$gt": 0}}
	update := Document{
		"$inc": Document{"price": 10},
		"$set": Document{"updated": true},
	}
	err := mockDocDB.UpdateMany("products", filter, update)
	assert.NoError(t, err)

	// Verify updates
	results, err := mockDocDB.FindDocument("products", nil)
	assert.NoError(t, err)
	for _, doc := range results {
		stock := int(doc["stock"].(float64))
		if stock > 0 {
			assert.Equal(t, true, doc["updated"])
			expectedPrice := int(doc["price"].(float64))
			if doc["name"] == "Product A" {
				assert.Equal(t, 110, expectedPrice)
			} else if doc["name"] == "Product B" {
				assert.Equal(t, 210, expectedPrice)
			}
		}
	}
}
