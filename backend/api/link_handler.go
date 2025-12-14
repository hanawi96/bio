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

func (h *LinkHandler) UpdateAllGroupStyles(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	
	if err := h.linkService.UpdateAllGroupStyles(userID, req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	return c.JSON(fiber.Map{"message": "All group styles updated successfully"})
}

func (h *LinkHandler) DuplicateLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	link, err := h.linkService.Duplicate(userID, linkID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *LinkHandler) BulkAction(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req struct {
		LinkIDs []string `json:"link_ids"`
		Action  string   `json:"action"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.linkService.BulkAction(userID, req.LinkIDs, req.Action); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(fiber.Map{"message": "Bulk action completed"})
}

func (h *LinkHandler) TogglePin(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	link, err := h.linkService.TogglePin(userID, linkID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(link)
}

func (h *LinkHandler) ReorderAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.linkService.ReorderWithBlocks(userID, req.Items); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"message": "Reordered successfully"})
}

// Group handlers
func (h *LinkHandler) CreateGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req struct {
		Title  string `json:"title"`
		Layout string `json:"layout"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	group, err := h.linkService.CreateGroup(userID, req.Title, req.Layout)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *LinkHandler) AddToGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	link, err := h.linkService.AddToGroup(userID, groupID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *LinkHandler) MoveToGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	var req struct {
		GroupID string `json:"group_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	link, err := h.linkService.MoveToGroup(userID, linkID, req.GroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(link)
}

func (h *LinkHandler) RemoveFromGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	link, err := h.linkService.RemoveFromGroup(userID, linkID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(link)
}

func (h *LinkHandler) DuplicateGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")
	group, err := h.linkService.DuplicateGroup(userID, groupID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *LinkHandler) ReorderGroupLinks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")
	var req struct {
		LinkIDs []string `json:"link_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	if err := h.linkService.ReorderGroupLinks(userID, groupID, req.LinkIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"message": "Group links reordered successfully"})
}
