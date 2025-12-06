package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/service"
)

type LinkHandler struct {
	linkService *service.LinkService
}

func NewLinkHandler(linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

func (h *LinkHandler) GetLinks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	links, err := h.linkService.GetByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Return empty array instead of null
	if links == nil {
		return c.JSON([]interface{}{})
	}
	return c.JSON(links)
}

func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	link, err := h.linkService.Create(userID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *LinkHandler) UpdateLink(c *fiber.Ctx) error {
	linkID := c.Params("id")
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	link, err := h.linkService.Update(linkID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(link)
}

func (h *LinkHandler) DeleteLink(c *fiber.Ctx) error {
	linkID := c.Params("id")
	if err := h.linkService.Delete(linkID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
