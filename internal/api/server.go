package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	server := &Server{db: db}
	router := gin.Default()

	// --- PUBLIC ROUTES ---
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", server.registerUser)
		authRoutes.POST("/login", server.loginUser)
	}

	// --- PROTECTED ROUTES ---
	// We create a new route group for the API that will use our middleware
	apiRoutes := router.Group("/api")
	apiRoutes.Use(authMiddleware()) // <-- Middleware is applied here
	{
		// Add a test route to verify that it works
		apiRoutes.GET("/me", server.getMe)

		// Service routes
		apiRoutes.POST("/services", server.createService)
		apiRoutes.GET("/services", server.getServices)
		// apiRoutes.GET("/services/:id", server.getService)
		// apiRoutes.PUT("/services/:id", server.updateService)
		apiRoutes.DELETE("/services/:id", server.deleteService)
	}

	server.router = router
	return server
}

// ... (función Start sin cambios) ...

// getMe es un handler de prueba para ver el ID del usuario autenticado.
func (s *Server) getMe(c *gin.Context) {
	// Obtenemos el userID que nuestro middleware añadió al contexto
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// El userID se guardó como int64, así que lo convertimos
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a protected route!",
		"user_id": userID.(int64),
	})
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
