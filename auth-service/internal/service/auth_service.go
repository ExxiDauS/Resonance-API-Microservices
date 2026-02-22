package service

import (
	"errors"

	"auth-service/internal/repository/postgres"
	"auth-service/internal/security"

	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db}
}

func (s *AuthService) Register(email, password, name string) error {
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}
	u := postgres.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		Username:     name,
	}
	if err := s.db.Create(&u).Error; err != nil {
		if strings.Contains(err.Error(), "users_email_key") ||
			strings.Contains(err.Error(), "email") {
			return errors.New("email already exists")
		}
		if strings.Contains(err.Error(), "users_username_key") ||
			strings.Contains(err.Error(), "username") {
			return errors.New("username already exists")
		}
		return err
	}
	return nil
}

func (s *AuthService) Login(identifier, password string) (*postgres.User, string, error) {

	var u postgres.User

	err := s.db.
		Where("email = ? OR username = ?", identifier, identifier).
		First(&u).Error

	if err != nil {
		return nil, "", errors.New("invalid email/username or password")
	}

	if !security.CheckPassword(u.PasswordHash, password) {
		return nil, "", errors.New("invalid email/username or password")
	}

	token, err := security.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	return &u, token, nil
}
