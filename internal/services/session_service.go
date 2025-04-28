package services

import (
	"auth-go/internal/database/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SessionService interface {
	IsSessionActive(token string) error
	Create(session models.Session) error
	FindByToken(token string) (models.Session, error)
	FindByUserID(userID string) ([]models.Session, error)
	RevokeByRefreshToken(token string) error
	RevokeAllByUserID(userID string) error
	ListByUserID(userID string) ([]models.Session, error)
}

type sessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) SessionService {
	return &sessionService{db: db}
}

func (s *sessionService) IsSessionActive(token string) error {
	var session models.Session
	err := s.db.Where("refresh_token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("session not active or not found")
	}
	return err
}

func (s *sessionService) Create(session models.Session) error {
	return s.db.Create(&session).Error
}

func (s *sessionService) FindByToken(token string) (models.Session, error) {
	var session models.Session
	err := s.db.Where("refresh_token = ?", token).First(&session).Error
	return session, err
}

func (s *sessionService) FindByUserID(userID string) ([]models.Session, error) {
	var sessions []models.Session
	err := s.db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

func (s *sessionService) RevokeByRefreshToken(token string) error {
	result := s.db.Model(&models.Session{}).Where("refresh_token = ? AND is_revoked = ?", token, false).Updates(map[string]interface{}{"is_revoked": true, "revoked_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no active session found for this refresh token")
	}
	return nil
}

func (s *sessionService) RevokeAllByUserID(userID string) error {
	return s.db.Model(&models.Session{}).Where("user_id = ? AND is_revoked = ?", userID, false).Updates(map[string]interface{}{"is_revoked": true, "revoked_at": time.Now()}).Error
}

func (s *sessionService) ListByUserID(userID string) ([]models.Session, error) {
	var sessions []models.Session
	err := s.db.Where("user_id = ?", userID).Find(&sessions).Error
	return sessions, err
}

func (s *sessionService) CreateSession(userID string, refreshToken string, deviceID string, expiresAt time.Time) (*models.Session, error) {
	session := models.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		DeviceID:     deviceID,
		ExpiresAt:    expiresAt,
		IsRevoked:    false,
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
