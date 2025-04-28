package services

import (
	"errors"
	"gorm.io/gorm"
	"auth-go/internal/database/models"
)

type UserService interface {
	AddUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id string) (models.User, error)
	IsEmailRegistered(email string) (bool, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) AddUser(user models.User) error {
	return s.db.Create(&user).Error
}

func (s *userService) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (s *userService) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := s.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (s *userService) IsEmailRegistered(email string) (bool, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // no existe
	}
	if err != nil {
		return false, err // error de DB
	}
	return true, nil // sí existe
}
