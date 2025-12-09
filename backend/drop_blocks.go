package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS blocks CASCADE;")
	if err != nil {
		log.Fatal("Failed to drop table:", err)
	}

	fmt.Println("✅ Dropped blocks table")

	// Re-run migration
	migrationSQL, err := os.ReadFile("migrations/004_create_blocks_table.sql")
	if err != nil {
		log.Fatal("Failed to read migration:", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	fmt.Println("✅ Re-created blocks table with new schema")
}
