package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/yourusername/linkbio/api"
	"github.com/yourusername/linkbio/config"
	"github.com/yourusername/linkbio/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize config
	cfg := config.New()

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migration for image_shape (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE links 
		ADD COLUMN IF NOT EXISTS image_shape VARCHAR(20) DEFAULT 'square'
	`)
	if err != nil {
		log.Println("⚠️ Migration warning:", err)
	} else {
		log.Println("✅ Migration: image_shape column ready")
	}

	// Update image_placement constraint to include 'alternating'
	_, err = db.Exec(`
		ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_image_placement;
		ALTER TABLE links ADD CONSTRAINT chk_image_placement 
		CHECK (image_placement IN ('left', 'right', 'top', 'bottom', 'alternating'))
	`)
	if err != nil {
		log.Println("⚠️ Image placement constraint warning:", err)
	} else {
		log.Println("✅ Migration: image_placement constraint updated")
	}

	// Add description column (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE links 
		ADD COLUMN IF NOT EXISTS description TEXT
	`)
	if err != nil {
		log.Println("⚠️ Description migration warning:", err)
	} else {
		log.Println("✅ Migration: description column ready")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "LinkBio API",
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(middleware.RateLimiter())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	apiGroup := app.Group("/api")
	api.SetupRoutes(apiGroup, db, cfg)

	// Start scheduler for auto-publish/unpublish
	scheduler := api.GetScheduler()
	if scheduler != nil {
		scheduler.Start()
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
