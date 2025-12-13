package api

import (
	"fmt"
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
	fmt.Printf("📤 UploadLinkThumbnail START: linkID=%s\n", linkID)

	// Get uploaded file
	file, err := c.FormFile("thumbnail")
	if err != nil {
		fmt.Printf("❌ No file uploaded: %v\n", err)
		return fiber.NewError(fiber.StatusBadRequest, "No file uploaded")
	}
	fmt.Printf("✅ File received: %s, Size: %d bytes, Type: %s\n", file.Filename, file.Size, file.Header.Get("Content-Type"))

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		fmt.Printf("❌ File too large: %d bytes\n", file.Size)
		return fiber.NewError(fiber.StatusBadRequest, "File size must be less than 5MB")
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" && contentType != "image/gif" {
		fmt.Printf("❌ Invalid file type: %s\n", contentType)
		return fiber.NewError(fiber.StatusBadRequest, "Only JPEG, PNG, WebP, and GIF images are allowed")
	}
	fmt.Printf("✅ File validation passed\n")

	// Open file
	src, err := file.Open()
	if err != nil {
		fmt.Printf("❌ Failed to open file: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Upload to Cloudinary
	fmt.Printf("📤 Uploading to Cloudinary...\n")
	thumbnailURL, err := utils.UploadToCloudinary(src, file.Filename)
	if err != nil {
		fmt.Printf("❌ Cloudinary upload failed: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload image: "+err.Error())
	}
	fmt.Printf("✅ Uploaded to Cloudinary: %s\n", thumbnailURL)

	// Update link with thumbnail URL
	data := map[string]interface{}{
		"thumbnail_url": thumbnailURL,
	}
	fmt.Printf("📝 Updating link with thumbnail URL...\n")
	link, err := h.linkService.Update(linkID, data)
	if err != nil {
		fmt.Printf("❌ Failed to update link: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update link")
	}
	fmt.Printf("✅ Link updated successfully\n")

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

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "No file uploaded")
	}

	// Validate file size (max 10MB for videos, 5MB for images)
	contentType := file.Header.Get("Content-Type")
	isVideo := contentType == "video/mp4" || contentType == "video/webm"
	maxSize := int64(5 * 1024 * 1024) // 5MB default
	if isVideo {
		maxSize = 10 * 1024 * 1024 // 10MB for videos
	}
	
	fmt.Printf("📁 File upload: name=%s, type=%s, size=%d bytes\n", file.Filename, contentType, file.Size)
	
	if file.Size > maxSize {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("File size must be less than %dMB", maxSize/(1024*1024)))
	}

	// Validate file type
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
		"video/mp4":  true,
		"video/webm": true,
	}
	
	if !validTypes[contentType] {
		return fiber.NewError(fiber.StatusBadRequest, "Only JPEG, PNG, WebP, GIF images and MP4, WebM videos are allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Upload to Cloudinary (use appropriate function based on file type)
	var fileURL string
	if isVideo {
		fileURL, err = utils.UploadVideoToCloudinary(src, file.Filename)
	} else {
		fileURL, err = utils.UploadToCloudinary(src, file.Filename)
	}
	
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload file: "+err.Error())
	}

	fmt.Printf("✅ Upload successful: %s\n", fileURL)
	return c.JSON(fiber.Map{
		"url": fileURL,
	})
}
