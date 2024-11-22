package api

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/klemis/user-actions-api/storage"
	"github.com/klemis/user-actions-api/types"
)

type Server struct {
	listenAddr string
	router     *gin.Engine
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		router:     gin.Default(),
		store:      store,
	}
}

func (s *Server) Start() error {
	s.router.GET("/users/:id", s.handleGetUserByID)
	s.router.GET("/users/referal-index", s.handleGetReferralIndex)
	s.router.GET("/users/:id/actions/count", s.handleGetActionCountByUserID)
	s.router.GET("/actions/:type/next-probalility", s.handleGetNextActionProbability)

	return s.router.Run(s.listenAddr)
}

// handleGetUserByID handles getting a user
func (s *Server) handleGetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve user data from the store.
	user := s.store.GetUser(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleGetActionCountByUserID handles getting the total number of actions for a given user ID.
func (s *Server) handleGetActionCountByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve action count.
	count := s.store.CountActionsByUserID(userID)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (s *Server) handleGetNextActionProbability(c *gin.Context) {
	actionType := c.Param("type")
	if actionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action type is required"})
		return
	}

	// Retrieve all actions sorted by user and createdAt.
	actions := s.store.GetActions()

	actionCounts := make(map[string]int)
	totalNextActions := 0

	// Count next actions after each specified action type.
	for i := 0; i < len(actions)-1; i++ {
		if actions[i].Type == actionType && actions[i].UserID == actions[i+1].UserID {
			nextAction := actions[i+1].Type
			actionCounts[nextAction]++
			totalNextActions++
		}
	}

	// Calculate probabilities.
	var result = make(types.ActionsProbalibity)
	for action, count := range actionCounts {
		probability := float64(count) / float64(totalNextActions)
		result[action] = math.Round(probability*100) / 100
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleGetReferralIndex(c *gin.Context) {
	// Retrieve all actions.
	actions := s.store.GetActions()
	if len(actions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No actions found"})
		return
	}

	// Create a mapping of users to the IDs of users they referred.
	referrals := make(types.Referral)
	for _, action := range actions {
		if action.Type == "REFER_USER" && action.TargetUser != 0 {
			referrals[action.UserID] = append(referrals[action.UserID], action.TargetUser)
		}
	}

	if len(referrals) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No referrals found"})
		return
	}

	// Calculate referral index for each user.
	referralIndex := make(types.ReferralIndex)
	for userId := range referrals {
		visited := make(map[int]bool)

		var dfs func(int)
		dfs = func(user int) {
			if visited[user] {
				return
			}

			visited[user] = true
			// Traverse each referral made by the current user.
			for _, referredUser := range referrals[user] {
				dfs(referredUser)
			}

			referralIndex[userId]++
		}
		// Start DFS on each referred user in the referrals list for userId.
		for _, referredUser := range referrals[userId] {
			dfs(referredUser)
		}
	}

	// TODO: display also users with 0 value?

	c.JSON(http.StatusOK, referralIndex)
}
