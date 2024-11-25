package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateAccessToken(string, string) (AuthResponse, error)
}

// JWTService generates access and refresh tokens
type jwtService struct {
	secret string
	accessTokenTTL time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(secret string, accessTokenTTL, refreshTokenTTL time.Duration) JWTService {
	return &jwtService{
		secret: secret,
		accessTokenTTL: accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// Represents custom claims using JWT
type UserClaim struct {
	UserID string `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Generates access token
func (s *jwtService) GenerateAccessToken(userID, role string) (AuthResponse, error) {
	now := time.Now()

	// Generate access token
	userClaims := &UserClaim{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
	accessToken, err := s.generateToken(userClaims)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &UserClaim{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := s.generateToken(refreshClaims)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return AuthResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Generates a JWT given the userClaim and signs it
func (j *jwtService) generateToken(claims *UserClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secret))
}