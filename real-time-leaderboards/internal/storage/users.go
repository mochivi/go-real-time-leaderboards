package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mochivi/go-real-time-leaderboards/internal/models"
)

type UserRepo interface {
	Create(context.Context, *models.RegisterUser, string) (*models.User, error)
	GetByUsername(context.Context, string) (*models.User, error)
	GetByID(context.Context, string) (*models.User, error)
	Update(context.Context, *models.UpdateUser) (*models.User, error)
	Delete(context.Context, string) error
}

// Users will be added to postgres users table
type UserRepoPG struct {
	db *sql.DB
}

func NewUserRepoPG(db *sql.DB) *UserRepoPG {
	return &UserRepoPG{
		db: db,
	}
}

func (ur *UserRepoPG) Create(ctx context.Context, registerUser *models.RegisterUser, passwordHash string) (*models.User, error) {

	stmt, err := ur.db.PrepareContext(ctx, `
		INSERT INTO users (username, password_hash, email)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, role, created_at, updated_at
	`)
	if err != nil {
		log.Printf("failed to prepare user creation statement: %v", err)
		return nil, fmt.Errorf("failed to prepare user creation statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var createdUser models.User
	if err = stmt.QueryRowContext(
		ctx,
		registerUser.Username,
		passwordHash,
		registerUser.Email,
	).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.Role,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	); err != nil {
		log.Printf("failed to execute user creation query: %v", err)
		return nil, fmt.Errorf("failed to execute user creation query: %w", err)
	}
	
	return &createdUser, nil
}


func (ur *UserRepoPG) GetByUsername(ctx context.Context, username string) (*models.User, error) {

	stmt, err := ur.db.PrepareContext(ctx, `
		SELECT id, username, password_hash, email, role, created_at, updated_at
		FROM users
		WHERE username = $1`,
	)
	if err != nil {
		log.Printf("failed to prepare query user by username: %v", err)
		return nil, fmt.Errorf("failed to prepare query user by username: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	if err = stmt.QueryRowContext(ctx, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		log.Printf("failed to query user by username: %v", err)
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}
	
	return &user, nil
}

func (ur *UserRepoPG) GetByID(ctx context.Context, userID string) (*models.User, error) {

	stmt, err := ur.db.PrepareContext(ctx, `
		SELECT id, username, password_hash, email, role, created_at, updated_at
		FROM users
		WHERE id = $1`,
	)
	if err != nil {
		log.Printf("failed to prepare query user by username: %v", err)
		return nil, fmt.Errorf("failed to prepare query user by username: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	if err = stmt.QueryRowContext(ctx, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		log.Printf("failed to query user by username: %v", err)
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}
	
	return &user, nil
}

func (ur *UserRepoPG) Update(ctx context.Context, updateUser *models.UpdateUser) (*models.User, error) {

	stmt, err := ur.db.PrepareContext(ctx, `
		UPDATE users 
		SET 
			username = $1,
			email = $2,
			role= $3
		WHERE id = $4
		RETURNING id, username, email, role, created_at, updated_at
	`)
	if err != nil {
		log.Printf("failed to prepare user creation statement: %v", err)
		return nil, fmt.Errorf("failed to prepare user creation statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var updatedUser models.User
	if err = stmt.QueryRowContext(
		ctx,
		updateUser.Username,
		updateUser.Email,
		updateUser.Role,
		updateUser.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.Role,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	); err != nil {
		log.Printf("failed to execute user creation query: %v", err)
		return nil, fmt.Errorf("failed to execute user creation query: %w", err)
	}
	
	return &updatedUser, nil
}

func (ur *UserRepoPG) Delete(ctx context.Context, userID string) error {

	stmt, err := ur.db.PrepareContext(ctx, `DELETE FROM users WHERE id = $1`)
	if err != nil {
		log.Printf("failed to prepare user creation statement: %v", err)
		return fmt.Errorf("failed to prepare user creation statement: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		log.Printf("Failed to delete user with of id '%s': %v", userID, err)
		return fmt.Errorf("failed to delete user with of id '%s': %w", userID, err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("Failed to delete user with of id '%s', no rows affected: %v", userID, err)
		return fmt.Errorf("failed to delete user with of id '%s', no rows affected: %w", userID, err)
	}

	return nil
}