package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/klemis/user-actions-api/types"
)

// Storage interface for accessing user and action data.
type Storage interface {
	GetUser(int) *types.User
}

// InMemoryStorage implements the Storage interface with in-memory data.
type InMemoryStorage struct {
	users   map[int]types.User
	actions map[int]types.Action
	mu      sync.RWMutex
}

// NewInMemoryStorage loads data from JSON files and initializes storage.
func NewInMemoryStorage(userFile, actionFile string) (Storage, error) {
	storage := &InMemoryStorage{
		users:   make(map[int]types.User),
		actions: make(map[int]types.Action),
	}

	if err := storage.loadUsers(userFile); err != nil {
		return nil, fmt.Errorf("failed to load users: %v", err)
	}
	if err := storage.loadActions(actionFile); err != nil {
		return nil, fmt.Errorf("failed to load actions: %v", err)
	}

	return storage, nil
}

// Get retrieves a user by ID.
func (s *InMemoryStorage) GetUser(id int) *types.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil
	}

	return &user
}

// loadUsers reads and parses users.json file.
func (s *InMemoryStorage) loadUsers(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var users []types.User
	if err := json.Unmarshal(data, &users); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range users {
		s.users[user.ID] = user
	}

	return nil
}

// loadActions reads and parses actions.json file.
func (s *InMemoryStorage) loadActions(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var actions []types.Action
	if err := json.Unmarshal(data, &actions); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, action := range actions {
		s.actions[action.ID] = action
	}

	return nil
}
