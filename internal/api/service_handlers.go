package api

import (
	"context"
	"net/http"
	"strconv"
	"uptime-monitor/internal/database/db"

	"github.com/gin-gonic/gin"
)

type serviceInput struct {
	Name                 string `json:"name" binding:"required"`
	Target               string `json:"target" binding:"required,url"`
	CheckIntervalSeconds int    `json:"check_interval_seconds" binding:"required,min=30"`
}

// createService creates a new service for the authenticated user.
func (s *Server) createService(c *gin.Context) {
	var input serviceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	userID := c.GetInt64("userID")

	params := db.CreateServiceParams{
		UserID:               userID,
		Name:                 input.Name,
		Target:               input.Target,
		CheckIntervalSeconds: int64(input.CheckIntervalSeconds),
	}

	service, err := s.q.CreateService(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// getServices retrieves all services for the authenticated user.
func (s *Server) getServices(c *gin.Context) {
	userID := c.GetInt64("userID")

	services, err := s.q.GetServicesForUser(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	// Return an empty slice if no services are found, instead of null
	if services == nil {
		services = []db.Service{}
	}

	c.JSON(http.StatusOK, services)
}

// deleteService deletes a service by ID for the authenticated user.
func (s *Server) deleteService(c *gin.Context) {
	idParam := c.Param("id")
	serviceID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	userID := c.GetInt64("userID")

	params := db.DeleteServiceParams{
		ID:     serviceID,
		UserID: userID,
	}

	rowsAffected, err := s.q.DeleteService(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found or you do not have permission to delete it"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// getServiceStatusHistory retrieves the latest status checks for a specific service.
func (s *Server) getServiceStatusHistory(c *gin.Context) {
	idParam := c.Param("id")
	serviceID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	userID := c.GetInt64("userID")

	params := db.GetStatusChecksForServiceParams{
		ServiceID: serviceID,
		UserID:    userID, // Ensures the user owns the service
	}

	statusChecks, err := s.q.GetStatusChecksForService(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve status history"})
		return
	}

	// If the query returns no rows, it might be because the service doesn't exist
	// or doesn't belong to the user. We return an empty slice for simplicity.
	if statusChecks == nil {
		statusChecks = []db.StatusCheck{}
	}

	c.JSON(http.StatusOK, statusChecks)
}
