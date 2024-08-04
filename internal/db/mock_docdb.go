package db

import (
	"errors"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/docdb"
)

type MockDocDB struct {
	clusters   map[string]interface{}
	instances  map[string]interface{}
	documents  map[string][]interface{}
	lock       sync.RWMutex
	mockConfig *MockConfig
}

// CreateInstance implements DocumentDB.
func (*MockDocDB) CreateInstance(input *docdb.CreateDBInstanceInput) (*docdb.CreateDBInstanceOutput, error) {
	panic("unimplemented")
}

// DeleteCluster implements DocumentDB.
func (*MockDocDB) DeleteCluster(input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error) {
	panic("unimplemented")
}

// DeleteDocument implements DocumentDB.
func (*MockDocDB) DeleteDocument(collection string, filter interface{}) error {
	panic("unimplemented")
}

// DeleteInstance implements DocumentDB.
func (*MockDocDB) DeleteInstance(input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error) {
	panic("unimplemented")
}

// InsertDocument implements DocumentDB.
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

type MockConfig struct {
	SimulateLatency bool
	LatencyMs       int
	ErrorMode       bool
}

func NewMockDocDB(config *MockConfig) *MockDocDB {
	return &MockDocDB{
		clusters:   make(map[string]interface{}),
		instances:  make(map[string]interface{}),
		documents:  make(map[string][]interface{}),
		mockConfig: config,
	}
}

func (m *MockDocDB) CreateCluster(input *docdb.CreateDBClusterInput) (*docdb.CreateDBClusterOutput, error) {
	if m.mockConfig.ErrorMode {
		return nil, errors.New("simulated error")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mockConfig.SimulateLatency {
		time.Sleep(time.Duration(m.mockConfig.LatencyMs) * time.Millisecond)
	}
	clusterID := *input.DBClusterIdentifier
	m.clusters[clusterID] = input
	return &docdb.CreateDBClusterOutput{
		DBCluster: &docdb.DBCluster{
			DBClusterIdentifier: &clusterID,
		},
	}, nil
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
			// Implement a simple filter logic: match by "name" field
			if docMap, ok := doc.(map[string]interface{}); ok {
				if filterMap, ok := filter.(map[string]interface{}); ok {
					if docMap["name"] == filterMap["name"] {
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
