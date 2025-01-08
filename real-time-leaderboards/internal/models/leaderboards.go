package models

import (
	"fmt"
	"time"
)

type LeaderboardRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Live        bool      `json:"live"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (l *LeaderboardRequest) AddUpdatedAt() {
	l.UpdatedAt = time.Now()
}

type UpdateLeaderboardRequest struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Live        bool      `json:"live"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Identify which fields changes have been submitted to
func (l *UpdateLeaderboardRequest) AddUpdatedAt() {
	l.UpdatedAt = time.Now()
}

type Leaderboard struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Live        bool               `json:"live"`
	Entries     []LeaderboardEntry `json:"entries"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func (l Leaderboard) RedisKey() string {
	return fmt.Sprintf("leaderboard:%s", l.ID)
}

type LeaderboardEntryRequest struct {
	LeaderboardID string    `json:"leaderboard_id"`
	UserID        string    `json:"user_id"`
	Score         int       `json:"score"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (l *LeaderboardEntryRequest) AddUpdatedAt() {
	l.UpdatedAt = time.Now()
}

type LeaderboardEntry struct {
	ID            string    `json:"id"`
	LeaderboardID string    `json:"leaderboard_id"`
	User          User      `json:"user"`
	Score         int       `json:"score"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
