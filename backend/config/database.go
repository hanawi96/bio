package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(cfg *Config) (*sql.DB, error) {
	println("[Database] Connecting to:", cfg.DatabaseURL)
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	println("[Database] Connected successfully!")
	
	// Test permissions
	var currentUser, currentDB string
	err = db.QueryRow("SELECT current_user, current_database()").Scan(&currentUser, &currentDB)
	if err == nil {
		println("[Database] Current user:", currentUser, "Database:", currentDB)
	} else {
		println("[Database] Error checking user:", err.Error())
	}

	return db, nil
}
