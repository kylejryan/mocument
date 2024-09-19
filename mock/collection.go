package mock

import (
	"errors"
	"fmt"
	"time"

	"github.com/kylejryan/mocument/internal/utils"
	"github.com/kylejryan/mocument/logger"
	"go.uber.org/zap"
)

func init() {
	logger.Init()
}

func (m *MockDocDB) InsertDocument(collection string, document Document) error {
	if m.mockConfig.ErrorMode {
		return errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	m.documents[collection] = append(m.documents[collection], document)
	return nil
}

func (m *MockDocDB) InsertMany(collection string, documents []interface{}) error {
	if m.mockConfig.ErrorMode {
		return errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	m.documents[collection] = append(m.documents[collection], documents...)
	return nil
}

func (m *MockDocDB) UpdateMany(collection string, filter, update interface{}) error {
	if m.mockConfig.ErrorMode {
		return errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	if documents, ok := m.documents[collection]; ok {
		for i, doc := range documents {
			if docMap, ok := doc.(map[string]interface{}); ok {
				if filterMap, ok := filter.(map[string]interface{}); ok {
					match := true
					for key, value := range filterMap {
						if docMap[key] != value {
							match = false
							break
						}
					}
					if match {
						for k, v := range update.(map[string]interface{}) {
							docMap[k] = v
						}
						documents[i] = docMap
					}
				}
			}
		}
		return nil
	}
	return errors.New("document not found")
}

func (m *MockDocDB) UpdateOne(collection string, filter, update interface{}) error {
	if m.mockConfig.ErrorMode {
		return errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	if documents, ok := m.documents[collection]; ok {
		for i, doc := range documents {
			if docMap, ok := doc.(map[string]interface{}); ok {
				if filterMap, ok := filter.(map[string]interface{}); ok {
					match := true
					for key, value := range filterMap {
						if docMap[key] != value {
							match = false
							break
						}
					}
					if match {
						for k, v := range update.(map[string]interface{}) {
							docMap[k] = v
						}
						documents[i] = docMap
						return nil
					}
				}
			}
		}
		return errors.New("document not found")
	}
	return errors.New("collection not found")
}

func (m *MockDocDB) FindDocument(collection string, filter Document) ([]Document, error) {
	if m.mockConfig.ErrorMode {
		logger.Get().Debug("Simulated error in FindDocument", zap.String("collection", collection))
		return nil, errors.New("simulated error")
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	// Simulate latency if enabled
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	documents, ok := m.documents[collection]
	if !ok {
		return nil, errors.New("collection not found")
	}
	var results []Document
	for _, doc := range documents {
		if utils.MatchesFilter(doc, filter) {
			results = append(results, doc)
		}
	}
	return results, nil
}

func (m *MockDocDB) DeleteDocument(collection string, filter interface{}) error {
	if m.mockConfig.ErrorMode {
		return errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	if documents, ok := m.documents[collection]; ok {
		for i, doc := range documents {
			if docMap, ok := doc.(map[string]interface{}); ok {
				if filterMap, ok := filter.(map[string]interface{}); ok {
					match := true
					for key, value := range filterMap {
						if docMap[key] != value {
							match = false
							break
						}
					}
					if match {
						fmt.Printf("Deleting document: %+v\n", docMap)
						// Delete the document by removing it from the slice
						m.documents[collection] = append(documents[:i], documents[i+1:]...)
						return nil
					}
				}
			}
		}
		fmt.Println("No matching document found for deletion.")
		return errors.New("no matching document found")
	}
	fmt.Println("Collection not found.")
	return errors.New("collection not found")
}

func (m *MockDocDB) DeleteMany(collection string, filter Document) (int, error) {
	if m.mockConfig.ErrorMode {
		logger.Error("Simulated error in DeleteMany", zap.String("collection", collection))
		return 0, errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	// Simulate latency if enabled
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	documents, ok := m.documents[collection]
	if !ok {
		return 0, errors.New("collection not found")
	}
	var newDocuments []Document
	deletedCount := 0
	for _, doc := range documents {
		if matching.MatchesFilter(doc, filter) {
			deletedCount++
			logger.Info("Deleting document", zap.Any("document", doc))
		} else {
			newDocuments = append(newDocuments, doc)
		}
	}
	m.documents[collection] = newDocuments
	return deletedCount, nil
}

func (m *MockDocDB) CountDocuments(collection string, filter interface{}) (int, error) {
	if m.mockConfig.ErrorMode {
		return 0, errors.New("simulated error")
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	if documents, ok := m.documents[collection]; ok {
		if filter == nil {
			return len(documents), nil
		}
		if filterMap, ok := filter.(map[string]interface{}); ok {
			count := 0
			for _, doc := range documents {
				if docMap, ok := doc.(map[string]interface{}); ok {
					match := true
					for key, value := range filterMap {
						if docMap[key] != value {
							match = false
							break
						}
					}
					if match {
						count++
					}
				}
			}
			return count, nil
		}
		return 0, errors.New("invalid filter format")
	}
	return 0, errors.New("collection not found")
}
