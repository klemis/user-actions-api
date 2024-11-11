package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/klemis/user-actions-api/storage"
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
	s.router.GET("/user/:id", s.handleGetUserByID)

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
