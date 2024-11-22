package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/klemis/user-actions-api/types"
)

// Storage interface for accessing user and action data.
type Storage interface {
	GetUser(int) *types.User
	CountActionsByUserID(userID int) int
	GetActions() []types.Action
}

// inMemoryStorage implements the Storage interface with in-memory data.
type inMemoryStorage struct {
	users   map[int]types.User
	actions []types.Action
	mu      sync.RWMutex
}

// NewInMemoryStorage loads data from JSON files and initializes storage.
func NewInMemoryStorage(userFile, actionFile string) (Storage, error) {
	storage := &inMemoryStorage{
		users:   make(map[int]types.User),
		actions: []types.Action{},
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
func (s *inMemoryStorage) GetUser(id int) *types.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil
	}

	// Return a copy of the user to prevent modification of the data.
	userCopy := user

	return &userCopy
}

// CountActionsByUserID returns the count of actions for a specific user ID.
func (s *inMemoryStorage) CountActionsByUserID(userID int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, action := range s.actions {
		if action.UserID == userID {
			count++
		}
	}

	return count
}

func (s *inMemoryStorage) GetActions() []types.Action {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy of the slice to prevent external modification.
	actionsCopy := make([]types.Action, len(s.actions))
	copy(actionsCopy, s.actions)

	return actionsCopy
}

// CreateAction inserts a new action into the actions slice while maintaining the sorted order.
// The function uses a binary search to determine the correct position for insertion.
// This ensures the actions slice remains sorted by UserID and CreatedAt.

// func (s *InMemoryStorage) CreateAction(action types.Action) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	// Find the appropriate index to insert the new action.
// 	idx := sort.Search(len(s.actions), func(i int) bool {
// 		if s.actions[i].UserID == action.UserID {
// 			return s.actions[i].CreatedAt.After(action.CreatedAt)
// 		}
// 		return s.actions[i].UserID > action.UserID
// 	})

// 	// Insert the new action while maintaining sorted order.
// 	s.actions = append(s.actions[:idx], append([]types.Action{action}, s.actions[idx:]...)...)
// }

// loadUsers reads and parses users.json file.
func (s *inMemoryStorage) loadUsers(filename string) error {
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
func (s *inMemoryStorage) loadActions(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var actions []types.Action
	if err := json.Unmarshal(data, &actions); err != nil {
		return err
	}

	// Sort actions by user and createdAt before storing them.
	sort.Slice(actions, func(i, j int) bool {
		if actions[i].UserID == actions[j].UserID {
			return actions[i].CreatedAt.Before(actions[j].CreatedAt)
		}
		return actions[i].UserID < actions[j].UserID
	})

	s.mu.Lock()
	defer s.mu.Unlock()
	s.actions = actions

	return nil
}
