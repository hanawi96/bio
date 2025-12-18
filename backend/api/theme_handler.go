package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/service"
)

type ThemeHandler struct {
	themeService *service.ThemeService
}

func NewThemeHandler(themeService *service.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}

// GetMyThemes retrieves all themes for the authenticated user
// GET /api/themes/my
func (h *ThemeHandler) GetMyThemes(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	themes, err := h.themeService.GetMyThemes(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve themes")
	}

	return c.JSON(themes)
}

// GetThemeByID retrieves a specific theme
// GET /api/themes/:id
func (h *ThemeHandler) GetThemeByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	themeID := c.Params("id")

	theme, err := h.themeService.GetThemeByID(themeID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(theme)
}

// CreateTheme creates a new theme
// POST /api/themes
func (h *ThemeHandler) CreateTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		Name        string                 `json:"name"`
		Description *string                `json:"description"`
		Config      map[string]interface{} `json:"config"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	theme, err := h.themeService.CreateTheme(userID, req.Name, req.Description, req.Config)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(theme)
}

// UpdateTheme updates a theme
// PUT /api/themes/:id
func (h *ThemeHandler) UpdateTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	themeID := c.Params("id")

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	theme, err := h.themeService.UpdateTheme(themeID, userID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(theme)
}

// DeleteTheme deletes a theme
// DELETE /api/themes/:id
func (h *ThemeHandler) DeleteTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	themeID := c.Params("id")

	if err := h.themeService.DeleteTheme(themeID, userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Theme deleted successfully",
	})
}

// PublishTheme makes a theme public
// POST /api/themes/:id/publish
func (h *ThemeHandler) PublishTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	themeID := c.Params("id")

	theme, err := h.themeService.PublishTheme(themeID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(theme)
}

// UnpublishTheme makes a theme private
// POST /api/themes/:id/unpublish
func (h *ThemeHandler) UnpublishTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	themeID := c.Params("id")

	theme, err := h.themeService.UnpublishTheme(themeID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(theme)
}

// GetPublicThemes retrieves public themes for marketplace
// GET /api/themes/public
func (h *ThemeHandler) GetPublicThemes(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	themes, err := h.themeService.GetPublicThemes(limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve public themes")
	}

	return c.JSON(themes)
}

// GetThemeBySlug retrieves a public theme by slug
// GET /api/themes/slug/:slug
func (h *ThemeHandler) GetThemeBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	theme, err := h.themeService.GetThemeBySlug(slug)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Theme not found")
	}

	return c.JSON(theme)
}

// ExportTheme exports a theme as JSON
// GET /api/themes/:id/export
func (h *ThemeHandler) ExportTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	themeID := c.Params("id")

	theme, err := h.themeService.ExportTheme(themeID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Return theme config as downloadable JSON
	c.Set("Content-Disposition", "attachment; filename="+theme.Name+".json")
	c.Set("Content-Type", "application/json")

	return c.JSON(fiber.Map{
		"name":        theme.Name,
		"description": theme.Description,
		"config":      theme.Config,
		"exported_at": theme.UpdatedAt,
	})
}

// ImportTheme imports a theme from JSON
// POST /api/themes/import
func (h *ThemeHandler) ImportTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		Name         string                 `json:"name"`
		Description  *string                `json:"description"`
		Config       map[string]interface{} `json:"config"`
		SourceThemeID *string               `json:"source_theme_id"` // Optional: if importing from marketplace
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	theme, err := h.themeService.ImportTheme(userID, req.Name, req.Description, req.Config, req.SourceThemeID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(theme)
}
