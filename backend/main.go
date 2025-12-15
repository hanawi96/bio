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

	// Add show_description column (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE links 
		ADD COLUMN IF NOT EXISTS show_description BOOLEAN DEFAULT true
	`)
	if err != nil {
		log.Println("⚠️ Show description migration warning:", err)
	} else {
		log.Println("✅ Migration: show_description column ready")
	}

	// Add link card background columns (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE links ADD COLUMN IF NOT EXISTS has_card_background BOOLEAN DEFAULT true NOT NULL;
		ALTER TABLE links ADD COLUMN IF NOT EXISTS card_background_color VARCHAR(7) DEFAULT '#ffffff';
		ALTER TABLE links ADD COLUMN IF NOT EXISTS card_background_opacity INT DEFAULT 100;
		ALTER TABLE links ADD COLUMN IF NOT EXISTS card_border_radius INT DEFAULT 12;
		ALTER TABLE links ADD COLUMN IF NOT EXISTS card_text_color VARCHAR(7) DEFAULT '#000000';
		ALTER TABLE links ADD COLUMN IF NOT EXISTS shadow_x INT DEFAULT 0;
		ALTER TABLE links ADD COLUMN IF NOT EXISTS shadow_y INT DEFAULT 4;
		ALTER TABLE links ADD COLUMN IF NOT EXISTS shadow_blur INT DEFAULT 10;
	`)
	if err != nil {
		log.Println("⚠️ Card background migration warning:", err)
	} else {
		log.Println("✅ Migration: card background columns ready")
	}
	
	// Add shadow constraints
	_, err = db.Exec(`
		ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_shadow_x;
		ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_shadow_y;
		ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_shadow_blur;
		ALTER TABLE links ADD CONSTRAINT chk_links_shadow_x CHECK (shadow_x >= -20 AND shadow_x <= 20);
		ALTER TABLE links ADD CONSTRAINT chk_links_shadow_y CHECK (shadow_y >= 0 AND shadow_y <= 20);
		ALTER TABLE links ADD CONSTRAINT chk_links_shadow_blur CHECK (shadow_blur >= 0 AND shadow_blur <= 40);
	`)
	if err != nil {
		log.Println("⚠️ Shadow constraints warning:", err)
	} else {
		log.Println("✅ Migration: shadow constraints ready")
	}

	// Add constraints for card background columns
	_, err = db.Exec(`
		ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_card_background_opacity;
		ALTER TABLE links DROP CONSTRAINT IF EXISTS chk_links_card_border_radius;
		ALTER TABLE links ADD CONSTRAINT chk_links_card_background_opacity 
		  CHECK (card_background_opacity >= 0 AND card_background_opacity <= 100);
		ALTER TABLE links ADD CONSTRAINT chk_links_card_border_radius 
		  CHECK (card_border_radius >= 0 AND card_border_radius <= 32);
	`)
	if err != nil {
		log.Println("⚠️ Card background constraints warning:", err)
	} else {
		log.Println("✅ Migration: card background constraints ready")
	}

	// Add theme_config column to profiles (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE profiles 
		ADD COLUMN IF NOT EXISTS theme_config JSONB DEFAULT NULL
	`)
	if err != nil {
		log.Println("⚠️ Theme config migration warning:", err)
	} else {
		log.Println("✅ Migration: theme_config column ready")
	}

	// Add header_config column to profiles (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE profiles 
		ADD COLUMN IF NOT EXISTS header_config JSONB DEFAULT '{"layout":"centered","coverType":"gradient","coverColor":"#6366f1","coverGradientFrom":"#8b5cf6","coverGradientTo":"#ec4899","coverHeight":140,"avatarSize":110,"avatarBorder":4,"avatarBorderColor":"#ffffff","showCover":true,"bioAlign":"center","bioSize":"md"}'::jsonb
	`)
	if err != nil {
		log.Println("⚠️ Header config migration warning:", err)
	} else {
		log.Println("✅ Migration: header_config column ready")
	}
	
	// Update existing header_config with old avatar sizes
	_, err = db.Exec(`
		UPDATE profiles 
		SET header_config = jsonb_set(
			jsonb_set(header_config, '{avatarSize}', '110'),
			'{coverHeight}', '140'
		)
		WHERE (header_config->>'avatarSize')::int < 100
	`)
	if err != nil {
		log.Println("⚠️ Header config update warning:", err)
	} else {
		log.Println("✅ Updated old header configs with larger avatar sizes")
	}

	// Add index for theme_config (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_profiles_theme_config 
		ON profiles USING GIN (theme_config)
	`)
	if err != nil {
		log.Println("⚠️ Theme config index warning:", err)
	} else {
		log.Println("✅ Migration: theme_config index ready")
	}

	// Add social_links column to profiles (safe - uses IF NOT EXISTS)
	_, err = db.Exec(`
		ALTER TABLE profiles 
		ADD COLUMN IF NOT EXISTS social_links TEXT
	`)
	if err != nil {
		log.Println("⚠️ Social links migration warning:", err)
	} else {
		log.Println("✅ Migration: social_links column ready")
	}

	// Page settings migration
	_, err = db.Exec(`
		ALTER TABLE profiles 
		ADD COLUMN IF NOT EXISTS show_share_button BOOLEAN DEFAULT true,
		ADD COLUMN IF NOT EXISTS show_subscribe_button BOOLEAN DEFAULT true,
		ADD COLUMN IF NOT EXISTS hide_branding BOOLEAN DEFAULT false
	`)
	if err != nil {
		log.Println("⚠️ Page settings migration warning:", err)
	} else {
		log.Println("✅ Migration: page settings columns ready")
	}

	// Theme system refactor migration - Migrate link styles to theme_config
	_, err = db.Exec(`
		UPDATE profiles p
		SET theme_config = COALESCE(p.theme_config, '{}'::jsonb) || jsonb_build_object(
			'textAlignment', COALESCE(
				(SELECT l.text_alignment FROM links l WHERE l.profile_id = p.id AND l.is_group = true LIMIT 1),
				COALESCE(p.theme_config->>'textAlignment', 'center')
			),
			'textSize', COALESCE(
				(SELECT l.text_size FROM links l WHERE l.profile_id = p.id AND l.is_group = true LIMIT 1),
				COALESCE(p.theme_config->>'textSize', 'M')
			),
			'imageShape', COALESCE(
				(SELECT l.image_shape FROM links l WHERE l.profile_id = p.id AND l.is_group = true LIMIT 1),
				COALESCE(p.theme_config->>'imageShape', 'square')
			)
		)
		WHERE EXISTS (SELECT 1 FROM links l WHERE l.profile_id = p.id AND l.is_group = true)
	`)
	if err != nil {
		log.Println("⚠️ Theme refactor migration warning:", err)
	} else {
		log.Println("✅ Migration: Link styles migrated to theme_config")
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
