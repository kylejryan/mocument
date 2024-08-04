// internal/handlers/handler_test.go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylejryan/mocument/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestInsertDocument(t *testing.T) {
	mockConfig := &db.MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := db.NewMockDocDB(mockConfig)
	handler := NewHandler(mockDocDB)

	doc := map[string]interface{}{"name": "test"}
	docBytes, _ := json.Marshal(doc)

	req, err := http.NewRequest("POST", "/insert", bytes.NewBuffer(docBytes))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler.insertDocument).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

/*
func TestFindDocument(t *testing.T) {
	mockConfig := &db.MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := db.NewMockDocDB(mockConfig)
	handler := NewHandler(mockDocDB)

	mockDocDB.InsertDocument("collection", map[string]interface{}{"name": "test"})

	req, err := http.NewRequest("GET", "/find?filter=name=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler.findDocument).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var result map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&result)
	assert.Equal(t, "test", result["name"])
}
*/
// Implement similar tests for updateDocuments and deleteDocument...
