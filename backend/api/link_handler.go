package api

import (
	"fmt"
	
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
	
	// Get query parameters
	search := c.Query("search", "")
	status := c.Query("status", "")
	layoutType := c.Query("layout_type", "")
	sortBy := c.Query("sort_by", "")
	
	// Debug log
	fmt.Printf("🔍 GetLinks - search: '%s', status: '%s', layoutType: '%s', sortBy: '%s'\n", search, status, layoutType, sortBy)
	
	links, err := h.linkService.GetByUserIDWithFilters(userID, search, status, layoutType, sortBy)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	fmt.Printf("✅ Found %d links\n", len(links))
	
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
