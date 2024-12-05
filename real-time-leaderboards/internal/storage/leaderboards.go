package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mochivi/real-time-leaderboards/internal/models"
)

type LeaderboardRepo interface {
	Get(context.Context, string) (*models.Leaderboard, error)
	Create(context.Context, *models.Leaderboard) (*models.Leaderboard, error)
	CreateEntry(context.Context, *models.LeaderboardEntry) (*models.LeaderboardEntry, error)
	Update(context.Context, *models.Leaderboard) (*models.Leaderboard, error)
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
	
	// Get leaderboard 	
	query := fmt.Sprintf(`
		SELECT
			l.id
			,l.name
			,l.description 
			,l.live
			,l.created_at
			,l.updated_at
			,e.id
			,e.user_id
			,e.score
			,e.created_at
			,e.updated_at
			,u.username
		FROM leaderboards l
		LEFT JOIN leaderboard_entries e ON leaderboards.id = leaderboard_entries.leaderboard_id
		LEFT JOIN users u ON leaderboard_entries.user_id = users.id
		WHERE id = %s`,
		leaderboardID,
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	
	rows, err := lr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard models.Leaderboard
	for rows.Next() {
		var entry models.LeaderboardEntry
		if err := rows.Scan(
			&leaderboard.ID,
			&leaderboard.Name,
			&leaderboard.Description,
			&leaderboard.Live,
			&leaderboard.CreatedAt,
			&leaderboard.UpdatedAt,
			&entry.ID,
			&entry.Score,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.User.Username,
		); err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return &leaderboard, nil
}

func (lr *LeaderboardRepoPG) Create(ctx context.Context, leaderboard *models.Leaderboard) (*models.Leaderboard, error) {

	stmt, err := lr.db.PrepareContext(
		ctx, 
		`INSERT INTO leaderboards (id, name, description, live)
			VALUES (%s, %s, %s, %s, %s, %s)
			RETURNING id, name, description, live, created_at, updated_at`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var returnLeaderboard models.Leaderboard
	if err := stmt.QueryRowContext(
		ctx,
		leaderboard.ID,
		leaderboard.Name,
		leaderboard.Description,
		leaderboard.Live,
	).Scan(
		&returnLeaderboard.ID,
		&returnLeaderboard.Name,
		&returnLeaderboard.Description,
		&returnLeaderboard.Live,
		&returnLeaderboard.CreatedAt,
		&returnLeaderboard.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create leaderboard: %w", err)
	}

	return nil, nil
}

func (lr *LeaderboardRepoPG) CreateEntry(ctx context.Context, entry *models.LeaderboardEntry) (*models.LeaderboardEntry, error) {

	stmt, err := lr.db.PrepareContext(ctx,`
		INSERT INTO leaderboard_entries (id, leaderboard_id, user_id, score)
		VALUES (?,?,?,?)
		RETURNING id, leaderboard_id, user_id, score, created_at, updated_at`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var returnEntry models.LeaderboardEntry
	if err := stmt.QueryRowContext(
		ctx,
		entry.ID,
		entry.LeaderboardID,
		entry.User.ID,
		entry.Score,
	).Scan(
		&returnEntry.ID,
		&returnEntry.LeaderboardID,
		&returnEntry.User.ID,
		&returnEntry.Score,
		&returnEntry.CreatedAt,
		&returnEntry.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create leaderboard entry: %w", err)
	}


	return &returnEntry, nil
}

// Could be used for leaderboard updates when done by an admin, for example massive removal of invalid entries
func (lr *LeaderboardRepoPG) Update(ctx context.Context, leaderboard *models.Leaderboard) (*models.Leaderboard, error) {
	
	stmt, err := lr.db.PrepareContext(ctx, `
		UPDATE leaderboards (name, description, live, updated_at)
		VALUES (?,?,?,?)
		WHERE id = ?
		RETURNING id, name, description, live, created_at, updated_at`,
	)
	if err != nil {
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
		return nil, fmt.Errorf("failed to update leaderboard entry: %w", err)
	}

	return &updatedEntry, nil
}

func (lr *LeaderboardRepoPG) Delete(ctx context.Context, leaderboardID string) error {
	stmt, err := lr.db.PrepareContext(ctx, `DELETE FROM leaderboards WHERE id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := stmt.ExecContext(ctx, leaderboardID); err != nil {
		return fmt.Errorf("failed to delete leaderboard: %w", err)
	}

	return nil
}

