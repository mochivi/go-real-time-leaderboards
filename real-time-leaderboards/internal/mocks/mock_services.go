package mocks

import (
	"context"
	"time"

	"github.com/mochivi/go-real-time-leaderboards/internal/auth"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) CreateAccessTokens(userID string, role string) (auth.AuthResponse, error) {
	args := m.Called()
	return args.Get(0).(auth.AuthResponse), args.Error(1)
}

func (m *MockJWTService) VerifyToken(tokenString string) (*auth.CustomClaims, error) {
	args := m.Called()
	return args.Get(0).(*auth.CustomClaims), args.Error(1)
}

func (m *MockJWTService) ParseTokenFromHeader(authHeader string) (string, bool) {
	args := m.Called()
	return args.String(0), args.Bool(1)
}

type MockRedisService struct {
	mock.Mock
}

func (m *MockRedisService) Set(ctx context.Context, key string, value any, exp time.Duration) error {
	args := m.Called(key, value, exp)
	return args.Error(0)
}

func (m *MockRedisService) Get(ctx context.Context, key string, target any) error {
	args := m.Called(key, target)
	return args.Error(0)
}

func (m *MockRedisService) JSONSet(ctx context.Context, key string, path string, value any, exp time.Duration) error {
	args := m.Called(key, path, value, exp)
	return args.Error(0)
}

func (m *MockRedisService) JSONGet(ctx context.Context, key string, path string, target any) error {
	args := m.Called(key, path, target)
	return args.Error(0)
}
