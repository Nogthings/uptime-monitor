package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"uptime-monitor/internal/models"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

func GetServicesToUpdate(ctx context.Context, db *pgxpool.Pool) ([]models.Service, error) {
	rows, err := db.Query(ctx, "SELECT id, user_id, name, target, check_interval_seconds, created_at, status, last_checked_at, latency_ms FROM services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.Target, &s.CheckIntervalSeconds, &s.CreatedAt, &s.Status, &s.LastCheckedAt, &s.LatencyMS); err != nil {
			// It's better to log the error and continue, but for now we'll return
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}

func UpdateServiceCheck(ctx context.Context, db *pgxpool.Pool, serviceID int64, status string, latency int, lastChecked time.Time) error {
	_, err := db.Exec(ctx, "UPDATE services SET status = $1, latency_ms = $2, last_checked_at = $3 WHERE id = $4", status, latency, lastChecked, serviceID)
	return err
}
