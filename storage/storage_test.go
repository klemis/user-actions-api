package storage

import (
	"sync"
	"testing"
	"time"

	"github.com/klemis/user-actions-api/types"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	mockTime, err := time.Parse(time.RFC3339, "2021-07-04T12:47:09.888Z")
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}

	tests := []struct {
		name     string
		userID   int
		users    map[int]types.User
		expected *types.User
	}{
		{
			name:   "User exists",
			userID: 2,
			users: map[int]types.User{
				1: {ID: 1, Name: "Tom", CreatedAt: mockTime.Add(1 * time.Hour)},
				2: {ID: 2, Name: "Alice", CreatedAt: mockTime},
			},
			expected: &types.User{ID: 2, Name: "Alice", CreatedAt: mockTime},
		},
		{
			name:     "User does not exist",
			userID:   2,
			users:    map[int]types.User{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &InMemoryStorage{
				users: tt.users,
				mu:    sync.RWMutex{},
			}

			result := storage.GetUser(tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountActionsByUserID(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		actions  map[int]types.Action
		expected int
	}{
		{
			name:   "Multiple actions for user",
			userID: 1,
			actions: map[int]types.Action{
				1: {ID: 1, UserID: 1, Type: "WELCOME"},
				2: {ID: 2, UserID: 1, Type: "CONNECT_CRM"},
				3: {ID: 3, UserID: 2, Type: "EDIT_CONTACT"},
			},
			expected: 2,
		},
		{
			name:   "No actions for user",
			userID: 3,
			actions: map[int]types.Action{
				1: {ID: 1, UserID: 1, Type: "ADD_CONTACT"},
				2: {ID: 2, UserID: 2, Type: "VIEW_CONTACTS"},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &InMemoryStorage{
				actions: tt.actions,
				mu:      sync.RWMutex{},
			}

			result := storage.CountActionsByUserID(tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetActions(t *testing.T) {
	mockTime, err := time.Parse(time.RFC3339, "2021-07-04T12:47:09.888Z")
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}

	tests := []struct {
		name     string
		actions  map[int]types.Action
		expected []types.Action
	}{
		{
			name: "Get sorted actions",
			actions: map[int]types.Action{
				1: {ID: 1, UserID: 1, Type: "WELCOME", CreatedAt: mockTime},
				2: {ID: 2, UserID: 1, Type: "EDIT_CONTACT", CreatedAt: mockTime.Add(3 * time.Hour)},
				3: {ID: 3, UserID: 1, Type: "CONNECT_CRM", CreatedAt: mockTime.Add(1 * time.Hour)},
			},
			expected: []types.Action{
				{ID: 1, UserID: 1, Type: "WELCOME", CreatedAt: mockTime},
				{ID: 3, UserID: 1, Type: "CONNECT_CRM", CreatedAt: mockTime.Add(1 * time.Hour)},
				{ID: 2, UserID: 1, Type: "EDIT_CONTACT", CreatedAt: mockTime.Add(3 * time.Hour)},
			},
		},
		{
			name:     "No actions",
			actions:  map[int]types.Action{},
			expected: []types.Action{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &InMemoryStorage{
				actions: tt.actions,
				mu:      sync.RWMutex{},
			}

			result := storage.GetActions()
			assert.Equal(t, tt.expected, result)
		})
	}
}
