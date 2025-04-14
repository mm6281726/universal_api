package storage

import (
	"errors"
	"sync"
	"universal_api/internal/models"
)

// Storage interface for storing API docs
type Storage interface {
	SaveAPIDoc(doc *models.APIDoc) error
	GetAPIDoc(id string) (*models.APIDoc, error)
	GetAllAPIDocs() ([]*models.APIDoc, error)
}

// MemoryStorage implements Storage using in-memory storage
type MemoryStorage struct {
	docs  map[string]*models.APIDoc
	mutex sync.RWMutex
}

// NewMemoryStorage creates a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		docs: make(map[string]*models.APIDoc),
	}
}

// SaveAPIDoc saves an API doc to memory
func (s *MemoryStorage) SaveAPIDoc(doc *models.APIDoc) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if doc.ID == "" {
		return errors.New("API doc ID cannot be empty")
	}

	s.docs[doc.ID] = doc
	return nil
}

// GetAPIDoc gets an API doc from memory
func (s *MemoryStorage) GetAPIDoc(id string) (*models.APIDoc, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	doc, ok := s.docs[id]
	if !ok {
		return nil, errors.New("API doc not found")
	}

	return doc, nil
}

// GetAllAPIDocs gets all API docs from memory
func (s *MemoryStorage) GetAllAPIDocs() ([]*models.APIDoc, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	docs := make([]*models.APIDoc, 0, len(s.docs))
	for _, doc := range s.docs {
		docs = append(docs, doc)
	}

	return docs, nil
}

// SQLiteStorage implements Storage using SQLite
// This is a placeholder for future implementation
type SQLiteStorage struct {
	// DB connection would go here
}

// NewSQLiteStorage creates a new SQLiteStorage
func NewSQLiteStorage() *SQLiteStorage {
	return &SQLiteStorage{}
}

// SaveAPIDoc saves an API doc to SQLite
func (s *SQLiteStorage) SaveAPIDoc(doc *models.APIDoc) error {
	// This would be implemented to save to SQLite
	return errors.New("SQLite storage not implemented yet")
}

// GetAPIDoc gets an API doc from SQLite
func (s *SQLiteStorage) GetAPIDoc(id string) (*models.APIDoc, error) {
	// This would be implemented to get from SQLite
	return nil, errors.New("SQLite storage not implemented yet")
}

// GetAllAPIDocs gets all API docs from SQLite
func (s *SQLiteStorage) GetAllAPIDocs() ([]*models.APIDoc, error) {
	// This would be implemented to get all from SQLite
	return nil, errors.New("SQLite storage not implemented yet")
}
