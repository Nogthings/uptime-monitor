package models

import "time"

type Service struct {
	ID                   int64     `json:"id"`
	UserID               int64     `json:"-"`
	Name                 string    `json:"name"`
	Target               string    `json:"target"`
	CheckIntervalSeconds int64     `json:"check_interval_seconds"`
	CreatedAt            time.Time `json:"created_at"`
	Status               string    `json:"status"`
	LastCheckedAt        time.Time `json:"last_checked_at"`
	LatencyMS            int       `json:"latency_ms"`
}
