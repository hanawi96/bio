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
	
	profile, err := h.profileService.GetByUsername(username)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}

	return c.JSON(profile)
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
