package mocks

import (
	"github.com/mochivi/go-real-time-leaderboards/internal/auth"
	"github.com/stretchr/testify/mock"
)


type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) CreateAccessTokens(userID string, role string) (auth.AuthResponse, error)  {
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