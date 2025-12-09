package api

import (
	"github.com/gofiber/fiber/v2"
)





func (h *LinkHandler) DuplicateLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	
	link, err := h.linkService.Duplicate(userID, linkID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *LinkHandler) BulkAction(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	var req struct {
		LinkIDs []string `json:"link_ids"`
		Action  string   `json:"action"` // "delete", "activate", "deactivate"
	}
	
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	
	if len(req.LinkIDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No links selected")
	}
	
	if err := h.linkService.BulkAction(userID, req.LinkIDs, req.Action); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LinkHandler) TogglePin(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	
	link, err := h.linkService.TogglePin(userID, linkID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

	return c.SendStatus(fiber.StatusNoContent)
}

// CreateGroup creates a new link group
func (h *LinkHandler) CreateGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	var req struct {
		Title  string `json:"title"`
		Layout string `json:"layout"` // "list", "grid", "carousel"
	}
	
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	
	if req.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Group title is required")
	}
	
	if req.Layout == "" {
		req.Layout = "list"
	}
	
	group, err := h.linkService.CreateGroup(userID, req.Title, req.Layout)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	return c.Status(fiber.StatusCreated).JSON(group)
}

// AddToGroup adds a link to an existing group
func (h *LinkHandler) AddToGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")
	
	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	
	if req["title"] == nil || req["url"] == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Title and URL are required")
	}
	
	link, err := h.linkService.AddToGroup(userID, groupID, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	
	return c.Status(fiber.StatusCreated).JSON(link)
}

// GetGroupItems retrieves all items in a group


// MoveToGroup moves an existing link into a group
func (h *LinkHandler) MoveToGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	
	var req struct {
		GroupID string `json:"group_id"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	
	if req.GroupID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Group ID is required")
	}
	
	link, err := h.linkService.MoveToGroup(userID, linkID, req.GroupID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	
	return c.JSON(link)
}

// RemoveFromGroup removes a link from its group
func (h *LinkHandler) RemoveFromGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	linkID := c.Params("id")
	
	link, err := h.linkService.RemoveFromGroup(userID, linkID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	
	return c.JSON(link)
}

// DuplicateGroup duplicates a group and all its children
func (h *LinkHandler) DuplicateGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")
	
	group, err := h.linkService.DuplicateGroup(userID, groupID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	return c.Status(fiber.StatusCreated).JSON(group)
}

// ReorderGroupLinks reorders links within a group
func (h *LinkHandler) ReorderGroupLinks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")
	
	var req struct {
		LinkIDs []string `json:"link_ids"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}
	
	if len(req.LinkIDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No links provided")
	}
	
	if err := h.linkService.ReorderGroupLinks(userID, groupID, req.LinkIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	
	return c.SendStatus(fiber.StatusNoContent)
}
