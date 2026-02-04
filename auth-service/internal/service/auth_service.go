package service

import (
	"errors"

	"auth-service/internal/repository/postgres"
	"auth-service/internal/security"

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

	var exists postgres.User
	err := s.db.Where("email = ?", email).First(&exists).Error
	if err == nil {
		return errors.New("email already exists")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	u := postgres.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		DisplayName:  name,
	}

	return s.db.Create(&u).Error
}

func (s *AuthService) Login(email, password string) (*postgres.User, string, error) {

	var u postgres.User
	err := s.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, "", errors.New("invalid email/password")
	}

	if !security.CheckPassword(u.PasswordHash, password) {
		return nil, "", errors.New("invalid email/password")
	}

	token, err := security.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, "", err
	}

	return &u, token, nil
}
