package models

import "time"

type Leaderboard struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Live bool `json:"live"`
	Entries []LeaderboardEntry `json:"entries"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LeaderboardEntry struct {
	ID string `json:"id"`
	LeaderboardID string `json:"leaderboard_id"`
	User User `json:"user"`
	Score int `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}