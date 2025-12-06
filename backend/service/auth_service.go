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

func (s *AuthService) Register(email, username, password string) (interface{}, string, error) {
	println("[AuthService] Starting registration for email:", email, "username:", username)
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		println("[AuthService] Error hashing password:", err.Error())
		return nil, "", err
	}
	println("[AuthService] Password hashed successfully")

	user, err := s.userRepo.Create(email, username, string(hashedPassword))
	if err != nil {
		println("[AuthService] Error creating user in repository:", err.Error())
		return nil, "", err
	}
	println("[AuthService] User created successfully with ID:", user.ID)

	token, err := s.generateToken(user.ID)
	if err != nil {
		println("[AuthService] Error generating token:", err.Error())
		return nil, "", err
	}
	println("[AuthService] Token generated successfully")

	return user, token, nil
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
