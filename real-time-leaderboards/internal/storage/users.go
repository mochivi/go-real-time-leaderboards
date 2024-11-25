package storage

import (
	"context"
	"database/sql"

	"github.com/mochivi/real-time-leaderboards/internal/models"
)

type UserRepo interface {
	Get(context.Context, string) (*models.User, error)
	Create(context.Context, *models.User) (*models.User, error)
	Update(context.Context, *models.User) (*models.User, error)
	Delete(context.Context, string) error
}

type UserRepoPG struct {
	db *sql.DB
}

func NewUserRepoPG(db *sql.DB) *UserRepoPG {
	return &UserRepoPG{
		db: db,
	}
}

func (us *UserRepoPG) Get(ctx context.Context, userID string) (*models.User, error) {return nil, nil}
func (us *UserRepoPG) Create(ctx context.Context, user *models.User) (*models.User, error) {return nil, nil}
func (us *UserRepoPG) Update(ctx context.Context, user *models.User) (*models.User, error) {return nil, nil}
func (us *UserRepoPG) Delete(ctx context.Context, userID string) error {return nil}
