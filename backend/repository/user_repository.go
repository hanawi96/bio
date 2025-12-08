package repository

import (
	"database/sql"
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
	println("[UserRepository] Creating user - email:", email, "username:", username)
	
	// Test connection first
	var testUser string
	err := r.db.QueryRow("SELECT current_user").Scan(&testUser)
	if err != nil {
		println("[UserRepository] ERROR checking current user:", err.Error())
	} else {
		println("[UserRepository] Current DB user:", testUser)
	}
	
	var user User
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, password_hash
	`
	println("[UserRepository] Executing INSERT query...")
	println("[UserRepository] Query:", query)
	println("[UserRepository] Params: email=", email, "username=", username)
	
	err = r.db.QueryRow(query, email, username, passwordHash).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
	)
	if err != nil {
		println("[UserRepository] ERROR inserting into users table:", err.Error())
		return nil, err
	}
	println("[UserRepository] User inserted successfully with ID:", user.ID)

	// Create profile for user
	println("[UserRepository] Creating profile for user...")
	_, err = r.db.Exec(`INSERT INTO profiles (user_id) VALUES ($1)`, user.ID)
	if err != nil {
		println("[UserRepository] ERROR inserting into profiles table:", err.Error())
		return nil, err
	}
	println("[UserRepository] Profile created successfully")

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
