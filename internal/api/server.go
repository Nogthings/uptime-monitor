package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server represents the API server.
type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
}

// NewServer creates a new API server instance.
func NewServer(db *pgxpool.Pool) *Server {
	server := &Server{db: db}
	router := gin.Default()

	// Define routes here, e.g.:
	// router.POST("/register", server.handleRegister)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	server.router = router
	return server
}

// Start runs the API server on the specified address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
