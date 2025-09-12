package models

import "time"

type Service struct {
	ID                   int64     `json:"id"`
	UserID               int64     `json:"-"`
	Name                 string    `json:"name"`
	Target               string    `json:"target"`
	CheckIntervalSeconds int64     `json:"check_interval_seconds"`
	CreatedAt            time.Time `json:"created_at"`
}
