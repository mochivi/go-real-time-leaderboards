package models

import "time"

type Leaderboard struct {
	ID        string             `json:"id"`
	Entries   []LeaderboardEntry `json:"entries"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type LeaderboardEntry struct {
	LeaderboardID string        `json:"leaderboard_id"`
	UserID        string        `json:"user_id"`
	Score         time.Duration `json:"score"`
	Live          bool          `json:"live"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
