package utils

import (
	"auth-go/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Pair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenActions interface {
	GenerateToken(userID string) (*Pair, error)
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateAccessToken(accessToken string) (*jwt.RegisteredClaims, error)
	ValidateRefreshToken(refreshToken string) (*jwt.RegisteredClaims, error)
	ValidateToken(anyToken string) (*jwt.RegisteredClaims, error)
}

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewManager() (*Manager, error) {
	cfg := config.LoadConfig()
	if cfg.AccessTokenSecret == "" || cfg.RefreshTokenSecret == "" {
		return nil, errors.New("token secrets must be configured")
	}

	return &Manager{
		accessSecret:  []byte(cfg.AccessTokenSecret),
		refreshSecret: []byte(cfg.RefreshTokenSecret),
		accessTTL:     15 * time.Minute,
		refreshTTL:    30 * 24 * time.Hour,
	}, nil
}

var _ TokenActions = (*Manager)(nil)

func (m *Manager) GenerateToken(userID string) (*Pair, error) {
	acc, err := m.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	ref, err := m.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	return &Pair{AccessToken: acc, RefreshToken: ref}, nil
}

func (m *Manager) GenerateAccessToken(userID string) (string, error) {
	return m.sign(userID, m.accessSecret, m.accessTTL)
}

func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	return m.sign(userID, m.refreshSecret, m.refreshTTL)
}

func (m *Manager) ValidateAccessToken(token string) (*jwt.RegisteredClaims, error) {
	return m.parse(token, m.accessSecret)
}

func (m *Manager) ValidateRefreshToken(token string) (*jwt.RegisteredClaims, error) {
	return m.parse(token, m.refreshSecret)
}

func (m *Manager) ValidateToken(token string) (*jwt.RegisteredClaims, error) {
	if claims, err := m.ValidateAccessToken(token); err == nil {
		return claims, nil
	}
	return m.ValidateRefreshToken(token)
}

// --- helpers ---

func (m *Manager) sign(userID string, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func (m *Manager) parse(tokenStr string, secret []byte) (*jwt.RegisteredClaims, error) {
	parser := jwt.Parser{
		ValidMethods: []string{jwt.SigningMethodHS256.Alg()},
	}
	tok, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !tok.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := tok.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("claims conversion failed")
	}
	return claims, nil
}
