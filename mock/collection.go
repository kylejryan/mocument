package mock

import (
	"errors"
	"fmt"
	"time"
)

func (m *MockDocDB) InsertDocument(collection string, document interface{}) error {
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

func (m *MockDocDB) FindDocument(collection string, filter interface{}) (interface{}, error) {
	if m.mockConfig.ErrorMode {
		return nil, errors.New("simulated error")
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	if documents, ok := m.documents[collection]; ok {
		if filter == nil {
			return documents, nil
		}
		if filterMap, ok := filter.(map[string]interface{}); ok {
			var results []interface{}
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
						results = append(results, doc)
					}
				}
			}
			return results, nil
		}
		return nil, errors.New("invalid filter format")
	}
	return []interface{}{}, nil // Return an empty slice if collection is not found
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

func (m *MockDocDB) DeleteMany(collection string, filter interface{}) (int, error) {
	if m.mockConfig.ErrorMode {
		return 0, errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	if documents, ok := m.documents[collection]; ok {
		var remainingDocuments []interface{}
		deletedCount := 0
		for _, doc := range documents {
			if docMap, ok := doc.(map[string]interface{}); ok {
				match := true
				if filterMap, ok := filter.(map[string]interface{}); ok {
					for key, value := range filterMap {
						if docMap[key] != value {
							match = false
							break
						}
					}
				}
				if match {
					deletedCount++
					continue
				}
			}
			remainingDocuments = append(remainingDocuments, doc)
		}
		m.documents[collection] = remainingDocuments
		return deletedCount, nil
	}
	return 0, errors.New("collection not found")
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
