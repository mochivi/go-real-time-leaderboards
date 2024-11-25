package storage

import (
	"context"
	"database/sql"

	"github.com/mochivi/real-time-leaderboards/internal/models"
)

type LeaderboardRepo interface {
	Get(context.Context, string) (*models.Leaderboard, error)
	Create(context.Context, *models.Leaderboard) (*models.Leaderboard, error)
	AddEntry(context.Context, *models.LeaderboardEntry) (*models.Leaderboard, error)
	Update(context.Context, *models.Leaderboard) (*models.Leaderboard, error)
	Delete(context.Context) error
}

// Postgres implementation
type LeaderboardRepoPG struct {
	db *sql.DB
}

func NewLeaderboardRepoPG(db *sql.DB) *LeaderboardRepoPG {
	return &LeaderboardRepoPG{
		db: db,
	}
}

// Interface implementations
func (lr *LeaderboardRepoPG) Get(ctx context.Context, leaderboardID string) (*models.Leaderboard, error) {
	return nil, nil
}

func (lr *LeaderboardRepoPG) Create(ctx context.Context, leaderboard *models.Leaderboard) (*models.Leaderboard, error) {
	return nil, nil
}

func (lr *LeaderboardRepoPG) AddEntry(ctx context.Context, entry *models.LeaderboardEntry) (*models.Leaderboard, error) {
	return nil, nil
}

// Could be used for leaderboard updates when done by an admin, for example massive removal of invalid entries
func (lr *LeaderboardRepoPG) Update(ctx context.Context, leaderboard *models.Leaderboard) (*models.Leaderboard, error) {
	return nil, nil
}

func (lr *LeaderboardRepoPG) Delete(ctx context.Context) error {
	return nil
}

