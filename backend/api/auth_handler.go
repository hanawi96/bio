package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/linkbio/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SetupUsernameRequest struct {
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		println("[AuthHandler] Error parsing request body:", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	println("[AuthHandler] Received registration request - email:", req.Email)

	user, token, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		println("[AuthHandler] Registration failed:", err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	println("[AuthHandler] Registration successful")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) SetupUsername(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	var req SetupUsernameRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.authService.SetupUsername(userID, req.Username); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *AuthHandler) CheckUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	available, err := h.authService.CheckUsernameAvailable(username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking username")
	}

	return c.JSON(fiber.Map{"available": available})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	return c.JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}
