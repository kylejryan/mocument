package main

import (
	"fmt"

	"github.com/kylejryan/mocument/mock"
)

func main() {
	mockConfig := &mock.MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDocDB := mock.NewMockDocDB(mockConfig)

	// Insert a document
	doc := map[string]interface{}{"name": "test"}
	err := mockDocDB.InsertDocument("collection", doc)
	if err != nil {
		fmt.Println("Error inserting document:", err)
		return
	}

	// Find the document
	filter := map[string]interface{}{"name": "test"}
	results, err := mockDocDB.FindDocument("collection", filter)
	if err != nil {
		fmt.Println("Error finding document:", err)
		return
	}

	fmt.Println("Found documents:", results)
}
