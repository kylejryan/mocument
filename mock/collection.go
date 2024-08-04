package mock

import (
	"errors"
	"time"
)

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
			if len(results) == 0 {
				return nil, errors.New("no matching document found")
			}
			return results, nil
		}
		return nil, errors.New("invalid filter format")
	}
	return nil, errors.New("document not found")
}

func (m *MockDocDB) DeleteDocument(collection string, filter interface{}) error {
	panic("unimplemented")
}
