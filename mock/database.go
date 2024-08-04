package mock

import (
	"errors"
	"sync"
	"time"
)

type Database struct {
	Name        string
	Collections map[string]*Collection
	lock        sync.RWMutex
	config      *MockConfig
}

type Collection struct {
	Name string
	Data []interface{}
}

func NewDatabase(name string, config *MockConfig) *Database {
	return &Database{
		Name:        name,
		Collections: make(map[string]*Collection),
		config:      config,
	}
}

func (db *Database) CreateCollection(name string) (*Collection, error) {
	if db.config.ErrorMode {
		return nil, errors.New("simulated error")
	}
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.config.SimulateLatency {
		time.Sleep(time.Duration(db.config.LatencyMs) * time.Millisecond)
	}
	collection := &Collection{Name: name}
	db.Collections[name] = collection
	return collection, nil
}

func (db *Database) GetCollection(name string) (*Collection, error) {
	if db.config.ErrorMode {
		return nil, errors.New("simulated error")
	}
	db.lock.RLock()
	defer db.lock.RUnlock()
	if db.config.SimulateLatency {
		time.Sleep(time.Duration(db.config.LatencyMs) * time.Millisecond)
	}
	collection, ok := db.Collections[name]
	if !ok {
		return nil, errors.New("collection not found")
	}
	return collection, nil
}

func (db *Database) DeleteCollection(name string) error {
	if db.config.ErrorMode {
		return errors.New("simulated error")
	}
	db.lock.Lock()
	defer db.lock.Unlock()
	if db.config.SimulateLatency {
		time.Sleep(time.Duration(db.config.LatencyMs) * time.Millisecond)
	}
	delete(db.Collections, name)
	return nil
}
