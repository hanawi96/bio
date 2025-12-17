package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/linkbio/config"
	"github.com/yourusername/linkbio/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Register(email, password string) (interface{}, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create temporary username unique from timestamp
	tempUsername := "temp_" + time.Now().Format("20060102150405")
	
	user, err := s.userRepo.Create(email, tempUsername, string(hashedPassword))
	if err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) SetupUsername(userID, username string) error {
	// Validate username
	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	// Check if username is available
	available, err := s.CheckUsernameAvailable(username)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("username already taken")
	}

	return s.userRepo.UpdateUsername(userID, username)
}

func (s *AuthService) CheckUsernameAvailable(username string) (bool, error) {
	_, err := s.userRepo.GetByUsername(username)
	if err != nil {
		// Username not found = available
		return true, nil
	}
	// Username found = not available
	return false, nil
}

func (s *AuthService) Login(email, password string) (interface{}, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
