package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       string         `gorm:"index;not null" json:"user_id"`
	RefreshToken string         `gorm:"uniqueIndex;size:255;not null" json:"-"`
	UserAgent    string         `gorm:"size:255;not null" json:"user_agent"`
	IPAddress    string         `gorm:"size:45;not null" json:"ip_address"`
	DeviceInfo   string         `gorm:"type:json" json:"device_info"`
	DeviceID     string         `gorm:"size:64;not null;index" json:"device_id"`
	LastUsedAt   time.Time      `json:"last_used_at"`
	ExpiresAt    time.Time      `gorm:"index;not null" json:"expires_at"`
	IsRevoked    bool           `gorm:"default:false" json:"is_revoked"`
	RevokedAt    *time.Time     `json:"revoked_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
