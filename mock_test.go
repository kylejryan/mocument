package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	// Verify the document is deleted
	results, err = mockDocDB.FindDocument("collection", map[string]interface{}{"name": "test"})
	assert.NoError(t, err)
	if results != nil {
		assert.Equal(t, 0, len(results.([]interface{})))
	} else {
		assert.Nil(t, results)
	}

}

func loadJSONFixture(filePath string, t *testing.T) map[string]interface{} {
	data, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	var doc map[string]interface{}
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
	filter := map[string]interface{}{"ID": "txn001"}
	results, err := mockDocDB.FindDocument("transactions", filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results.([]interface{})))
	assert.Equal(t, "txn001", results.([]interface{})[0].(map[string]interface{})["ID"])
	assert.Equal(t, 150.75, results.([]interface{})[0].(map[string]interface{})["Amount"])
	assert.Equal(t, "USD", results.([]interface{})[0].(map[string]interface{})["Currency"])
	assert.Equal(t, "Pending", results.([]interface{})[0].(map[string]interface{})["Status"])

	// Verify the items in the transaction
	items := results.([]interface{})[0].(map[string]interface{})["Items"].([]interface{})
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "prod001", items[0].(map[string]interface{})["ProductID"])
	assert.Equal(t, 50.25, items[0].(map[string]interface{})["Price"])
	assert.Equal(t, 2.0, items[1].(map[string]interface{})["Quantity"])
}
