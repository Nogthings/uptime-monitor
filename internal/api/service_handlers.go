package api

import (
	"context"
	"net/http"
	"strconv"
	"uptime-monitor/internal/models"

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

	service := models.Service{
		UserID:               userID,
		Name:                 input.Name,
		Target:               input.Target,
		CheckIntervalSeconds: int64(input.CheckIntervalSeconds),
	}

	query := `INSERT INTO services (user_id, name, target, check_interval_seconds)
             VALUES ($1, $2, $3, $4)
             RETURNING id, created_at`
	err := s.db.QueryRow(context.Background(), query, service.UserID, service.Name, service.Target, service.CheckIntervalSeconds).
		Scan(&service.ID, &service.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// getServices retrieves all services for the authenticated user.
func (s *Server) getServices(c *gin.Context) {
	userID := c.GetInt64("userID")

	query := `SELECT id, name, target, check_interval_seconds, created_at
             FROM services WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.Query(context.Background(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}
	defer rows.Close()

	services := []models.Service{}
	for rows.Next() {
		var service models.Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Target, &service.CheckIntervalSeconds, &service.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan service data"})
			return
		}
		services = append(services, service)
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

	// Ensure the service belongs to the user
	query := `DELETE FROM services WHERE id = $1 AND user_id = $2`
	cmdTag, err := s.db.Exec(context.Background(), query, serviceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	// Check if any rows were affected
	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found or you do not have permission to delete it"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
