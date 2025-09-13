package monitoring

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"uptime-monitor/internal/config"
	"uptime-monitor/internal/database/db"
	"uptime-monitor/internal/notifications"

	"github.com/jackc/pgx/v5/pgtype"
)

// Monitor holds the dependencies for the monitoring worker.
type Monitor struct {
	q         *db.Queries
	notifier  *notifications.EmailNotifier
	lastCheck map[int64]time.Time // In-memory cache to respect check intervals
}

// NewMonitor creates a new Monitor instance.
func NewMonitor(cfg *config.Config, q *db.Queries) *Monitor {
	return &Monitor{
		q:         q,
		notifier:  notifications.NewEmailNotifier(cfg),
		lastCheck: make(map[int64]time.Time),
	}
}

// Start begins the monitoring loop.
func (m *Monitor) Start() {
	log.Println("Monitoring worker started")
	ticker := time.NewTicker(15 * time.Second) // Ticker runs more frequently
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkAllServices()
		}
	}
}

func (m *Monitor) checkAllServices() {
	services, err := m.q.GetServicesAndOwners(context.Background())
	if err != nil {
		log.Printf("Error fetching services and owners: %v", err)
		return
	}

	for _, service := range services {
		// Respect the user-defined check interval
		if time.Since(m.lastCheck[service.ID]) > time.Duration(service.CheckIntervalSeconds)*time.Second {
			go m.checkService(service)
			m.lastCheck[service.ID] = time.Now()
		}
	}
}

func (m *Monitor) checkService(s db.GetServicesAndOwnersRow) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	currentStatus := ""
	params := db.CreateStatusCheckParams{
		ServiceID: s.ID,
	}

	startTime := time.Now()
	resp, err := client.Get(s.Target)
	responseTime := time.Since(startTime)

	if err != nil {
		currentStatus = "down"
		params.Status = currentStatus
		params.ErrorMessage = pgtype.Text{String: err.Error(), Valid: true}
	} else {
		defer resp.Body.Close()
		params.StatusCode = pgtype.Int4{Int32: int32(resp.StatusCode), Valid: true}
		params.ResponseTimeMs = pgtype.Int4{Int32: int32(responseTime.Milliseconds()), Valid: true}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			currentStatus = "up"
		} else {
			currentStatus = "down"
			params.ErrorMessage = pgtype.Text{String: fmt.Sprintf("Non-2xx status code: %d", resp.StatusCode), Valid: true}
		}
		params.Status = currentStatus
	}

	// --- State Change Detection & Notification ---
	previousStatus, err := m.q.GetLatestStatusCheckForService(context.Background(), s.ID)
	// sql.ErrNoRows is okay, means it's the first check ever.
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ERROR: Could not get previous status for service %d: %v", s.ID, err)
		return // Don't save the new check if we can't verify the old one
	}

	// If status has changed, send notification
	if previousStatus != currentStatus && err != sql.ErrNoRows {
		log.Printf("STATE CHANGE for %s: %s -> %s. Sending notification.", s.Name, previousStatus, currentStatus)
		subject := fmt.Sprintf("Uptime Alert: %s is %s", s.Name, strings.ToUpper(currentStatus))
		body := fmt.Sprintf("Your service '%s' (%s) is now %s.\n\nChecked at: %s", s.Name, s.Target, currentStatus, time.Now().Format(time.RFC1123))
		if err := m.notifier.SendNotification(s.OwnerEmail, subject, body); err != nil {
			log.Printf("ERROR: Failed to send notification for service %d: %v", s.ID, err)
		}
	}

	// --- Save the current check to the database ---
	_, dbErr := m.q.CreateStatusCheck(context.Background(), params)
	if dbErr != nil {
		log.Printf("ERROR: Failed to save status check for service %d: %v", s.ID, dbErr)
	}
}
