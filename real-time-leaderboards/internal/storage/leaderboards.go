package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mochivi/real-time-leaderboards/internal/models"
)

type LeaderboardRepo interface {
	Get(context.Context, string) (*models.Leaderboard, error)
	GetEntries(context.Context, string) ([]models.LeaderboardEntry, error)
	Create(context.Context, *models.LeaderboardRequest) (*models.Leaderboard, error)
	CreateEntry(context.Context, *models.LeaderboardEntryRequest) (*models.LeaderboardEntry, error)
	Update(context.Context, *models.UpdateLeaderboardRequest) (*models.Leaderboard, error)
	Delete(context.Context, string) error
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
	log.Printf("Getting leaderboard %s from DB", leaderboardID)

	// Get leaderboard
	stmt, err := lr.db.PrepareContext(
		ctx,
		`SELECT
			id
			,name
			,description 
			,live
			,created_at
			,updated_at
		FROM leaderboards
		WHERE id = $1`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare get statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var leaderboard models.Leaderboard
	if err := stmt.QueryRowContext(ctx, leaderboardID).Scan(
		&leaderboard.ID,
		&leaderboard.Name,
		&leaderboard.Description,
		&leaderboard.Live,
		&leaderboard.CreatedAt,
		&leaderboard.UpdatedAt,
	); err != nil {
		log.Printf("failed to get leaderboard: %v", err)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return &leaderboard, nil
}

func (lr *LeaderboardRepoPG) GetEntries(ctx context.Context, leaderboardID string) ([]models.LeaderboardEntry, error) {
	log.Printf("Getting leaderboard %s from DB", leaderboardID)

	// Get leaderboard 	
	stmt, err := lr.db.PrepareContext(
		ctx,
		`SELECT
			e.id
			,e.score
			,e.created_at
			,e.updated_at
			,u.id
			,u.username
		FROM leaderboard_entries e
		LEFT JOIN users u 
			ON e.user_id = u.id 
		WHERE e.leaderboard_id = $1`)
	if err != nil {
		log.Printf("Failed to prepare get statement: %v", err)
		return nil, fmt.Errorf("failed to prepare get statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := stmt.QueryContext(ctx, leaderboardID)
	if err != nil {
		log.Printf("failed to get leaderboard entries: %v", err)
		return nil, fmt.Errorf("failed to get leaderboard entries: %w", err)
	}

	entries := make([]models.LeaderboardEntry, 0)
	for rows.Next() {
		var entry models.LeaderboardEntry
		if err = rows.Scan(
			&entry.ID,
			&entry.Score,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.User.ID,
			&entry.User.Username,
		); err != nil {
			log.Printf("failed to scan leaderboard entry: %v", err)
			return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		log.Printf("failed to scan leaderboard entries: %v", err)
		return nil, fmt.Errorf("failed to scan leaderboard entries: %w", err)
	}

	return entries, nil
}

func (lr *LeaderboardRepoPG) Create(ctx context.Context, newLeaderboard *models.LeaderboardRequest) (*models.Leaderboard, error) {

	stmt, err := lr.db.PrepareContext(
		ctx, 
		`INSERT INTO public.leaderboards (name, description, live, updated_At)
			VALUES ($1, $2, $3, $4)
			RETURNING id, name, description, live, created_at, updated_at`,
	)
	if err != nil {
		log.Printf("Failed to prepare insert statement: %v", err)
		return nil, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var returnLeaderboard models.Leaderboard
	if err := stmt.QueryRowContext(
		ctx,
		newLeaderboard.Name,
		newLeaderboard.Description,
		newLeaderboard.Live,
		newLeaderboard.UpdatedAt,
	).Scan(
		&returnLeaderboard.ID,
		&returnLeaderboard.Name,
		&returnLeaderboard.Description,
		&returnLeaderboard.Live,
		&returnLeaderboard.CreatedAt,
		&returnLeaderboard.UpdatedAt,
	); err != nil {
		log.Printf("Failed to execute leaderboard creation query: %v", err)
		return nil, fmt.Errorf("failed to create leaderboard: %w", err)
	}

	return &returnLeaderboard, nil
}

func (lr *LeaderboardRepoPG) CreateEntry(ctx context.Context, entry *models.LeaderboardEntryRequest) (*models.LeaderboardEntry, error) {

	stmt, err := lr.db.PrepareContext(ctx,`
		INSERT INTO leaderboard_entries (leaderboard_id, user_id, score, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, leaderboard_id, user_id, score, created_at, updated_at`,
	)
	if err != nil {
		log.Printf("Failed to prepare entry creation statement: %v", err)
		return nil, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var returnEntry models.LeaderboardEntry
	if err := stmt.QueryRowContext(
		ctx,
		entry.LeaderboardID,
		entry.UserID,
		entry.Score,
		entry.UpdatedAt,
	).Scan(
		&returnEntry.ID,
		&returnEntry.LeaderboardID,
		&returnEntry.User.ID,
		&returnEntry.Score,
		&returnEntry.CreatedAt,
		&returnEntry.UpdatedAt,
	); err != nil {
		log.Printf("Failed to insert leaderboard entry: %v", err)
		return nil, fmt.Errorf("failed to create leaderboard entry: %w", err)
	}


	return &returnEntry, nil
}

// Could be used for leaderboard updates when done by an admin, for example massive removal of invalid entries
func (lr *LeaderboardRepoPG) Update(ctx context.Context, leaderboard *models.UpdateLeaderboardRequest) (*models.Leaderboard, error) {
	
	stmt, err := lr.db.PrepareContext(ctx, `
		UPDATE leaderboards
		SET
			name = $1,
			description = $2,
			live = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, name, description, live, created_at, updated_at`,
	)
	if err != nil {
		log.Printf("Failed to prepare update leaderboard statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()
	var updatedLeaderboard models.Leaderboard

	if err := stmt.QueryRowContext(
		ctx,
		leaderboard.Name,
		leaderboard.Description,
		leaderboard.Live,
		leaderboard.UpdatedAt,
		leaderboard.ID,
	).Scan(
		&updatedLeaderboard.ID,
		&updatedLeaderboard.Name,
		&updatedLeaderboard.Description,
		&updatedLeaderboard.Live,
		&updatedLeaderboard.CreatedAt,
		&updatedLeaderboard.UpdatedAt,
	); err != nil {
		log.Printf("Failed to update leaderboard: %v", err)
		return nil, fmt.Errorf("failed to update leaderboard: %w", err)
	}

	return &updatedLeaderboard, nil
}

func (lr *LeaderboardRepoPG) UpdateEntry(ctx context.Context, leaderboardEntry *models.LeaderboardEntry) (*models.LeaderboardEntry, error) {
	stmt, err := lr.db.PrepareContext(ctx, `
		UPDATE leaderboards_entries (score, updated_at)
		VALUES (?,?)
		WHERE id = ? AND leaderboard_id = ?
		RETURNING id, leaderboard_id, user_id, score, created_at, updated_at`,
	)
	if err != nil {
		log.Printf("Failed to prepare entry update statement: %v", err)
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var updatedEntry models.LeaderboardEntry

	if err := stmt.QueryRowContext(
		ctx,
		leaderboardEntry.Score,
		leaderboardEntry.UpdatedAt,
		leaderboardEntry.ID,
		leaderboardEntry.LeaderboardID,
	).Scan(
		&updatedEntry.ID,
		&updatedEntry.LeaderboardID,
		&updatedEntry.User.ID,
		&updatedEntry.Score,
		&updatedEntry.CreatedAt,
		&updatedEntry.UpdatedAt,
	); err != nil {
		log.Printf("Failed to execute update entry query: %v", err)
		return nil, fmt.Errorf("failed to update leaderboard entry: %w", err)
	}

	return &updatedEntry, nil
}

func (lr *LeaderboardRepoPG) Delete(ctx context.Context, leaderboardID string) error {
	stmt, err := lr.db.PrepareContext(ctx, `DELETE FROM leaderboards WHERE id = $1`)
	if err != nil {
		log.Printf("Failed to prepare leaderboard delete statement: %v", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := stmt.ExecContext(ctx, leaderboardID); err != nil {
		log.Printf("Failed to execute leaderboard delete query: %v", err)
		return fmt.Errorf("failed to delete leaderboard: %w", err)
	}

	return nil
}

