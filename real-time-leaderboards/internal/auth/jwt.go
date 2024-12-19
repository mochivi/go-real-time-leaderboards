package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {

	// Generate AccessToken and RefreshToken given the UserID and Role
	CreateAccessTokens(string, string) (AuthResponse, error)

	// VerifyToken parses the token and returns the user claims
	VerifyToken(string) (*CustomClaims, error)

	// ParseTokenFromHeader will return a token from an authorization header
	// panics if the Authorization header is malformed
	ParseTokenFromHeader(string) (string, bool)
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
type CustomClaims struct {
	UserID string `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Generates access token
func (s *jwtService) CreateAccessTokens(userID, role string) (AuthResponse, error) {
	now := time.Now()

	// Generate access token
	customClaims := &CustomClaims{
		UserID: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}
	accessToken, err := s.generateToken(customClaims)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &CustomClaims{
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

// Parses bearer token from header but does not validate it
func (j *jwtService) ParseTokenFromHeader(authHeader string) (string, bool) {

	// Comes in Authorization Bearer format:
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2  || parts[0] != "Bearer" {
		return "", false
	}

	return parts[1], true
}

func (j *jwtService) VerifyToken(tokenString string) (*CustomClaims, error) {
	
	// Parse token with custom claims
	var claims *CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != "HS512" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	// Verify token validity
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	// Get custom claims from token
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("failed to get custom claims from token")
	}

	return claims, nil
}

// Generates a JWT given the CustomClaims and signs it
func (j *jwtService) generateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secret))
}