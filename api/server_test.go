package api

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/klemis/user-actions-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage mocks the InMemoryStorage for testing.
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetUser(id int) *types.User {
	args := m.Called(id)
	if user := args.Get(0); user != nil {
		return user.(*types.User)
	}
	return nil
}

// CountActionsByUserID is a mocked method that counts actions for a specific user ID.
func (m *MockStorage) CountActionsByUserID(userID int) int {
	args := m.Called(userID)
	return args.Int(0)
}

// GetActions is a mocked method that retrieves all actions.
func (m *MockStorage) GetActions() []types.Action {
	args := m.Called()
	if actions := args.Get(0); actions != nil {
		return actions.([]types.Action)
	}
	return nil
}

// TestHandleGetUserByID tests the handleGetUserByID endpoint.
func TestHandleGetUserByID(t *testing.T) {
	// Set up mock storage.
	mockStore := &MockStorage{}
	server := &Server{store: mockStore}

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/user/:id", server.handleGetUserByID)

	mockTime, err := time.Parse(time.RFC3339, "2021-07-04T12:47:09.888Z")
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}

	tests := []struct {
		name           string
		userID         string
		mockReturn     *types.User
		mockReturnErr  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid User ID",
			userID:         "1",
			mockReturn:     &types.User{ID: 2, Name: "Alice", CreatedAt: mockTime},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": 2, "name": "Alice", "createdAt": "2021-07-04T12:47:09.888Z"}`,
		},
		{
			name:           "Invalid User ID (non-numeric)",
			userID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid user ID"}`,
		},
		{
			name:           "User Not Found",
			userID:         "55",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "User not found"}`,
		},
	}

	// Loop through each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := strconv.Atoi(tt.userID)
			if err == nil {
				mockStore.On("GetUser", id).Return(tt.mockReturn)
			}

			// Create a request and response recorder.
			req, _ := http.NewRequest("GET", "/user/"+tt.userID, nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			// Check the response status code.
			assert.Equal(t, tt.expectedStatus, response.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
			}
			assert.JSONEq(t, tt.expectedBody, response.Body.String())
		})
	}
}

// TestHandleGetActionCountByUserID tests the handleGetActionCountByUserID endpoint.
func TestHandleGetActionCountByUserID(t *testing.T) {
	// Set up mock storage.
	mockStore := &MockStorage{}
	server := &Server{store: mockStore}

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/user/:id/actions/count", server.handleGetActionCountByUserID)

	tests := []struct {
		name           string
		userID         string
		mockReturn     int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid User ID with actions",
			userID:         "1",
			mockReturn:     5,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"count": 5}`,
		},
		{
			name:           "Valid User ID with no actions",
			userID:         "2",
			mockReturn:     0,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"count": 0}`,
		},
		{
			name:           "Invalid User ID (non-numeric)",
			userID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Invalid user ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := strconv.Atoi(tt.userID)
			if err == nil {
				mockStore.On("CountActionsByUserID", id).Return(tt.mockReturn)
			}

			req, _ := http.NewRequest("GET", "/user/"+tt.userID+"/actions/count", nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedStatus, response.Code)

			assert.JSONEq(t, tt.expectedBody, response.Body.String())
		})
	}
}

// TestHandleGetNextActionProbability tests the handleGetNextActionProbability endpoint.
func TestHandleGetNextActionProbability(t *testing.T) {
	// Set up mock storage.
	mockStore := &MockStorage{}
	server := &Server{store: mockStore}

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/actions/:type/next-probability", server.handleGetNextActionProbability)

	// Example actions in the storage.
	actions := []types.Action{
		{ID: 1, UserID: 1, Type: "WELCOME"},
		{ID: 2, UserID: 1, Type: "CONNECT_CRM"},
		{ID: 3, UserID: 1, Type: "ADD_CONTACT"},
		{ID: 4, UserID: 2, Type: "EDIT_CONTACT"},
		{ID: 5, UserID: 3, Type: "WELCOME"},
		{ID: 6, UserID: 3, Type: "VIEW_CONTACTS"},
	}

	tests := []struct {
		name           string
		actionType     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Probability after WELCOME action",
			actionType:     "WELCOME",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"CONNECT_CRM":0.5, "VIEW_CONTACTS":0.5}`,
		},
		{
			name:           "Probability after CONNECT_CRM action",
			actionType:     "CONNECT_CRM",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ADD_CONTACT": 1.00}`,
		},
		{
			name:           "Probability after non-existent action",
			actionType:     "UNKNOWN_ACTION",
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`,
		},
		{
			name:           "Missing action type",
			actionType:     "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Action type is required"}`,
		},
	}

	// Loop through each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore.On("GetActions").Return(actions)

			req, _ := http.NewRequest("GET", "/actions/"+tt.actionType+"/next-probability", nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedStatus, response.Code)

			assert.JSONEq(t, tt.expectedBody, response.Body.String())
		})
	}
}

// TestHandleGetReferralIndex tests the handleGetReferralIndex endpoint.
func TestHandleGetReferralIndex(t *testing.T) {
	tests := []struct {
		name           string
		mockActions    []types.Action
		expectedStatus int
		expectedBody   string
	}{

		{
			name:           "No actions",
			mockActions:    []types.Action{},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "No actions found"}`,
		},
		{
			name: "No referrals",
			mockActions: []types.Action{
				{ID: 1, UserID: 1, Type: "WELCOME", TargetUser: 2},
				{ID: 2, UserID: 2, Type: "ADD_CONTACT", TargetUser: 3},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error": "No referrals found"}`,
		},
		{
			name: "Referral index calculation",
			mockActions: []types.Action{
				{ID: 1, UserID: 1, Type: "REFER_USER", TargetUser: 2},
				{ID: 2, UserID: 2, Type: "REFER_USER", TargetUser: 3},
				{ID: 3, UserID: 3, Type: "REFER_USER", TargetUser: 4},
				{ID: 4, UserID: 1, Type: "REFER_USER", TargetUser: 5},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"1": 4, "2": 2, "3": 1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockStorage{}
			server := &Server{store: mockStore}

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.GET("/users/referal-index", server.handleGetReferralIndex)

			mockStore.On("GetActions").Return(tt.mockActions)

			req, _ := http.NewRequest("GET", "/users/referal-index", nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			assert.Equal(t, tt.expectedStatus, response.Code)

			assert.JSONEq(t, tt.expectedBody, response.Body.String())
		})
	}
}
