package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/service"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetPublicProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	
	data, err := h.profileService.GetPublicProfileWithLinks(username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}

	return c.JSON(data)
}

func (h *ProfileHandler) GetMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	profile, err := h.profileService.GetByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}

	return c.JSON(profile)
}

func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	profile, err := h.profileService.Update(userID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(profile)
}

// ApplyTheme applies a theme preset to profile and all groups
func (h *ProfileHandler) ApplyTheme(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req struct {
		ThemeConfig  map[string]interface{} `json:"theme_config"`
		CardStyles   map[string]interface{} `json:"card_styles"`
		TextStyles   string                 `json:"text_styles"`
		HeaderConfig map[string]interface{} `json:"header_config"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.profileService.ApplyTheme(userID, req.ThemeConfig, req.CardStyles, req.TextStyles, req.HeaderConfig)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(result)
}
