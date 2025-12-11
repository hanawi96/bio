package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/service"
)

type BlockHandler struct {
	service *service.BlockService
}

func NewBlockHandler(service *service.BlockService) *BlockHandler {
	return &BlockHandler{service: service}
}

func (h *BlockHandler) GetBlocks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	println("🔍 GetBlocks for userID:", userID)

	blocks, err := h.service.GetBlocks(userID)
	if err != nil {
		println("❌ GetBlocks error:", err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch blocks", "details": err.Error()})
	}

	println("✅ Found", len(blocks), "blocks")
	return c.JSON(blocks)
}

func (h *BlockHandler) CreateBlock(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	block, err := h.service.CreateBlock(userID, data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create block", "details": err.Error()})
	}

	return c.Status(201).JSON(block)
}

func (h *BlockHandler) UpdateBlock(c *fiber.Ctx) error {
	blockID := c.Params("id")

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		println("❌ Parse error:", err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	println("📝 Updating block", blockID, "with data:", data)

	block, err := h.service.UpdateBlock(blockID, data)
	if err != nil {
		println("❌ Update block error:", err.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update block", "details": err.Error()})
	}

	return c.JSON(block)
}

func (h *BlockHandler) DeleteBlock(c *fiber.Ctx) error {
	blockID := c.Params("id")

	if err := h.service.DeleteBlock(blockID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete block"})
	}

	return c.SendStatus(204)
}

func (h *BlockHandler) ReorderBlocks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var data struct {
		BlockIDs []string `json:"block_ids"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.ReorderBlocks(userID, data.BlockIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reorder blocks"})
	}

	return c.SendStatus(204)
}

func (h *BlockHandler) BulkDelete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var data struct {
		BlockIDs []string `json:"block_ids"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.BulkDeleteBlocks(userID, data.BlockIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete blocks"})
	}

	return c.SendStatus(204)
}

// ReorderGroupBlocks reorders blocks within a group
func (h *BlockHandler) ReorderGroupBlocks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	groupID := c.Params("groupId")

	var req struct {
		BlockIDs []string `json:"block_ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	if len(req.BlockIDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No blocks provided")
	}

	if err := h.service.ReorderGroupBlocks(userID, groupID, req.BlockIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
