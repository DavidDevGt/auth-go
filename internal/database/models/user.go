package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string         `gorm:"size:60;not null" json:"name"`
	Email           string         `gorm:"size:120;uniqueIndex;not null" json:"email"`
	PasswordHash    string         `gorm:"size:255;not null" json:"-"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Sessions        []Session      `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}
