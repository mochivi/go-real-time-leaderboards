package mocks

import (
	"context"

	"github.com/mochivi/go-real-time-leaderboards/internal/models"
	"github.com/stretchr/testify/mock"
)


type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, registerUser *models.RegisterUser, passwordHash string) (*models.User, error) {
	args := m.Called(registerUser)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, updateUser *models.UpdateUser) (*models.User, error) {
	args := m.Called(updateUser)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) Delete(ctx context.Context, userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}


type MockLeaderboardsRepo struct {
	mock.Mock
}

func (m *MockLeaderboardsRepo) Get(ctx context.Context, leaderboardID string) (*models.Leaderboard, error) {
	args := m.Called(leaderboardID)
	return args.Get(0).(*models.Leaderboard), args.Error(1)
}

func (m *MockLeaderboardsRepo) GetEntries(ctx context.Context, leaderboardID string) ([]models.LeaderboardEntry, error) {
	args := m.Called(leaderboardID)
	return args.Get(0).([]models.LeaderboardEntry), args.Error(1)
}

func (m *MockLeaderboardsRepo) Create(ctx context.Context, newLeaderboard *models.LeaderboardRequest) (*models.Leaderboard, error) {
	args := m.Called(newLeaderboard)
	return args.Get(0).(*models.Leaderboard), args.Error(1)
}

func (m *MockLeaderboardsRepo) CreateEntry(ctx context.Context, entry *models.LeaderboardEntryRequest) (*models.LeaderboardEntry, error) {
	args := m.Called(entry)
	return args.Get(0).(*models.LeaderboardEntry), args.Error(1)
}

func (m *MockLeaderboardsRepo) Update(ctx context.Context, leaderboard *models.UpdateLeaderboardRequest) (*models.Leaderboard, error) {
	args := m.Called(leaderboard)
	return args.Get(0).(*models.Leaderboard), args.Error(1)
}

func (m *MockLeaderboardsRepo) Delete(ctx context.Context, leaderboardID string) error {
	args := m.Called(leaderboardID)
	return args.Error(0)
}
