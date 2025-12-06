package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *LinkHandler) ReorderLinks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	var req struct {
		LinkIDs []string `json:"link_ids"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	if err := h.linkService.ReorderLinks(userID, req.LinkIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LinkHandler) GetAnalytics(c *fiber.Ctx) error {
	_ = c.Locals("userID").(string)
	
	// TODO: Implement analytics
	return c.JSON(fiber.Map{
		"total_clicks": 0,
		"total_views":  0,
		"links":        []interface{}{},
	})
}
