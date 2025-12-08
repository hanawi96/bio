package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/pkg/utils"
	"github.com/yourusername/linkbio/service"
)

type UploadHandler struct {
	linkService    *service.LinkService
	profileService *service.ProfileService
}

func NewUploadHandler(linkService *service.LinkService, profileService *service.ProfileService) *UploadHandler {
	return &UploadHandler{
		linkService:    linkService,
		profileService: profileService,
	}
}

func (h *UploadHandler) UploadLinkThumbnail(c *fiber.Ctx) error {
	linkID := c.Params("id")

	// Get uploaded file
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "No file uploaded")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return fiber.NewError(fiber.StatusBadRequest, "File size must be less than 5MB")
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" && contentType != "image/gif" {
		return fiber.NewError(fiber.StatusBadRequest, "Only JPEG, PNG, WebP, and GIF images are allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Upload to Cloudinary
	thumbnailURL, err := utils.UploadToCloudinary(src, file.Filename)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload image: "+err.Error())
	}

	// Update link with thumbnail URL
	data := map[string]interface{}{
		"thumbnail_url": thumbnailURL,
	}
	link, err := h.linkService.Update(linkID, data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update link")
	}

	return c.JSON(link)
}

func (h *UploadHandler) DeleteLinkThumbnail(c *fiber.Ctx) error {
	linkID := c.Params("id")

	// Update link to remove thumbnail
	data := map[string]interface{}{
		"thumbnail_url": nil,
	}
	link, err := h.linkService.Update(linkID, data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete thumbnail")
	}

	return c.JSON(link)
}

func (h *UploadHandler) UploadAvatar(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Get uploaded file
	file, err := c.FormFile("avatar")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "No file uploaded")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return fiber.NewError(fiber.StatusBadRequest, "File size must be less than 5MB")
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" && contentType != "image/gif" {
		return fiber.NewError(fiber.StatusBadRequest, "Only JPEG, PNG, WebP, and GIF images are allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Upload to Cloudinary
	avatarURL, err := utils.UploadToCloudinary(src, file.Filename)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload image: "+err.Error())
	}

	// Update profile with new avatar
	profile, err := h.profileService.Update(userID, map[string]interface{}{
		"avatar_url": avatarURL,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update profile")
	}

	return c.JSON(profile)
}
