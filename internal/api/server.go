package api

import (
	"net/http"
	"uptime-monitor/internal/database/db"
	"uptime-monitor/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server now also holds web handlers
type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
	q      *db.Queries
}

func NewServer(dbPool *pgxpool.Pool) *Server {
	server := &Server{
		db: dbPool,
		q:  db.New(dbPool),
	}
	router := gin.Default()
	server.router = router

	// Pass the server instance to the web handlers
	webHandlers := &web.Server{Q: server.q}

	// --- STATIC FILES ---
	router.StaticFS("/static", http.Dir("public"))

	// --- WEB ROUTES ---
	router.GET("/login", webHandlers.ShowLoginPage)
	router.POST("/login", webHandlers.PostLoginPage)
	router.GET("/logout", webHandlers.Logout)

	// Authenticated web routes
	dashboardGroup := router.Group("/")
	dashboardGroup.Use(web.AuthMiddleware())
	{
		dashboardGroup.GET("/dashboard", webHandlers.ShowDashboardPage)
	}

	// --- PUBLIC API ROUTES ---
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	authAPIRoutes := router.Group("/auth")
	{
		authAPIRoutes.POST("/register", server.registerUser)
		authAPIRoutes.POST("/login", server.loginUser)
	}

	// --- PROTECTED API ROUTES ---
	apiRoutes := router.Group("/api")
	apiRoutes.Use(AuthMiddleware()) // Note: This is the API middleware
	{
		apiRoutes.GET("/me", server.getMe)
		apiRoutes.POST("/services", server.createService)
		apiRoutes.GET("/services", server.getServices)
		apiRoutes.DELETE("/services/:id", server.deleteService)
		apiRoutes.GET("/services/:id/status", server.getServiceStatusHistory)
	}

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// getMe remains a method of the original Server struct for the API
func (s *Server) getMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a protected route!",
		"user_id": userID.(int64),
	})
}
