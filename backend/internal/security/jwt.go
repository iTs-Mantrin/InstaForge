package security

import (
	"fmt"
	"time"

	"instaforge/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims.
type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Role   string   `json:"role"`
	Scopes []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT creation and validation.
type JWTManager struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:        []byte(cfg.Secret),
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
	}
}

// GenerateToken creates a signed JWT for a user.
func (m *JWTManager) GenerateToken(userID, email, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.accessExpiry)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "instaforge",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken creates a long-lived refresh token.
func (m *JWTManager) GenerateRefreshToken(userID, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.refreshExpiry)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "instaforge-refresh",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates and parses a JWT token.
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// RefreshToken validates a refresh token and returns a new access token.
func (m *JWTManager) RefreshToken(refreshToken string) (string, time.Time, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	return m.GenerateToken(claims.UserID, claims.Email, claims.Role)
}
