package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/docdb"
)

type MockConfig struct {
	SimulateLatency bool
	LatencyMs       int
	ErrorMode       bool
}

type MockDocDB struct {
	clusters   map[string]interface{}
	instances  map[string]interface{}
	documents  map[string][]interface{}
	lock       sync.RWMutex
	mockConfig *MockConfig
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

func (m *MockDocDB) DeleteCluster(input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error) {
	panic("unimplemented")
}

func (m *MockDocDB) CreateInstance(input *docdb.CreateDBInstanceInput) (*docdb.CreateDBInstanceOutput, error) {
	panic("unimplemented")
}

func (m *MockDocDB) DeleteInstance(input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error) {
	panic("unimplemented")
}

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
