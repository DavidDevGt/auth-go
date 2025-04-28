package services

import (
	"auth-go/internal/database/models"
	"auth-go/internal/utils"
	"auth-go/pkg/validators"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password, userAgent, ip string) (*utils.Pair, error)
	Refresh(refreshToken string) (*utils.Pair, error)
	Logout(refreshToken string) error
	RevokeSession(refreshToken string) error
	Register(user models.User) error
}

type authService struct {
	db           *gorm.DB
	users        UserService
	sessions     SessionService
	tokenManager utils.TokenActions
}

func NewAuthService(db *gorm.DB, users UserService, sessions SessionService, tokenManager utils.TokenActions) AuthService {
	return &authService{
		db:           db,
		users:        users,
		sessions:     sessions,
		tokenManager: tokenManager,
	}
}

func (s *authService) Login(email, password, userAgent, ip string) (*utils.Pair, error) {
	if err := validators.ValidateLogin(email, password); err != nil {
		return nil, err
	}
	user, err := s.users.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}
	pair, err := s.tokenManager.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	session := models.Session{
		UserID:       user.ID,
		UserAgent:    userAgent,
		IPAddress:    ip,
		RefreshToken: pair.RefreshToken,
		DeviceInfo:   "{}", // <-- JSON válido vacío
		LastUsedAt:   time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.sessions.Create(session); err != nil {
		return nil, err
	}
	return pair, nil
}

func (s *authService) Refresh(refreshToken string) (*utils.Pair, error) {
	if err := validators.ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}
	if err := s.sessions.IsSessionActive(refreshToken); err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	claims, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	pair, err := s.tokenManager.GenerateToken(claims.Subject)
	if err != nil {
		return nil, err
	}
	_ = s.sessions.RevokeByRefreshToken(refreshToken)

	session := models.Session{
		UserID:       claims.Subject,
		UserAgent:    "",
		IPAddress:    "",
		RefreshToken: pair.RefreshToken,
		DeviceInfo:   "{}", // <-- JSON válido vacío
		LastUsedAt:   time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	_ = s.sessions.Create(session)
	return pair, nil
}

func (s *authService) Logout(refreshToken string) error {
	if err := validators.ValidateRefreshToken(refreshToken); err != nil {
		return err
	}
	return s.sessions.RevokeByRefreshToken(refreshToken)
}

func (s *authService) RevokeSession(refreshToken string) error {
	if err := validators.ValidateRefreshToken(refreshToken); err != nil {
		return err
	}
	return s.sessions.RevokeByRefreshToken(refreshToken)
}

func (s *authService) Register(user models.User) error {
	if err := validators.ValidateRegister(user.Name, user.Email, user.PasswordHash); err != nil {
		return err
	}
	exists, err := s.users.IsEmailRegistered(user.Email)
	if err != nil {
		return err // error de DB
	}
	if exists {
		return errors.New("email already registered")
	}
	user.PasswordHash, _ = utils.HashPassword(user.PasswordHash)
	return s.users.AddUser(user)
}
