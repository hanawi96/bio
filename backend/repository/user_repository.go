package repository

import (
	"database/sql"
	"errors"
	"strings"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, username, passwordHash string) (*User, error) {
	var user User
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, password_hash
	`
	
	err := r.db.QueryRow(query, email, username, passwordHash).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
	)
	if err != nil {
		// Check for duplicate key errors
		if strings.Contains(err.Error(), "users_email_key") {
			return nil, errors.New("email already exists")
		}
		if strings.Contains(err.Error(), "users_username_key") {
			return nil, errors.New("username already taken")
		}
		return nil, err
	}

	// Create profile for user
	_, err = r.db.Exec(`INSERT INTO profiles (user_id) VALUES ($1)`, user.ID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, username, password_hash FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*User, error) {
	var user User
	query := `SELECT id, email, username, password_hash FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	var user User
	query := `SELECT id, email, username, password_hash FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUsername(userID, username string) error {
	query := `UPDATE users SET username = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, username, userID)
	return err
}
