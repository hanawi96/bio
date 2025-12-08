package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/config"
	"github.com/yourusername/linkbio/middleware"
	"github.com/yourusername/linkbio/repository"
	"github.com/yourusername/linkbio/service"
)

func SetupRoutes(api fiber.Router, db *sql.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	linkRepo := repository.NewLinkRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	profileService := service.NewProfileService(profileRepo, userRepo, linkRepo)
	linkService := service.NewLinkService(linkRepo)

	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	profileHandler := NewProfileHandler(profileService)
	linkHandler := NewLinkHandler(linkService)
	uploadHandler := NewUploadHandler(linkService, profileService)

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Public profile view
	api.Get("/p/:username", profileHandler.GetPublicProfile)

	// Protected routes
	protected := api.Group("", middleware.AuthRequired(cfg))
	
	// Profile management
	protected.Get("/profile", profileHandler.GetMyProfile)
	protected.Put("/profile", profileHandler.UpdateProfile)

	// Link management
	protected.Get("/links", linkHandler.GetLinks)
	protected.Post("/links", linkHandler.CreateLink)
	protected.Put("/links/reorder", linkHandler.ReorderLinks)
	protected.Post("/links/bulk", linkHandler.BulkAction)
	protected.Post("/links/:id/duplicate", linkHandler.DuplicateLink)
	protected.Post("/links/:id/pin", linkHandler.TogglePin)
	protected.Put("/links/:id", linkHandler.UpdateLink)
	protected.Delete("/links/:id", linkHandler.DeleteLink)

	// Upload management
	protected.Post("/links/:id/thumbnail", uploadHandler.UploadLinkThumbnail)
	protected.Delete("/links/:id/thumbnail", uploadHandler.DeleteLinkThumbnail)
	protected.Post("/profile/avatar", uploadHandler.UploadAvatar)

	// Analytics
	protected.Get("/analytics", linkHandler.GetAnalytics)
}
