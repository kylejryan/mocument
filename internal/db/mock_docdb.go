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

// Going to implement what I wrote below later ...
/*
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
		// Implement filter logic here...
		return documents, nil
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
		// Implement update logic here...
		for i := range documents {
			// Update each document based on filter and update...
		}
		return nil
	}
	return errors.New("document not found")
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
		// Implement delete logic here...
		return nil
	}
	return errors.New("document not found")
}
*/
