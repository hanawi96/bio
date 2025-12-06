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
	Username string `json:"username"`
	Password string `json:"password"`
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
	println("[AuthHandler] Received registration request - email:", req.Email, "username:", req.Username)

	user, token, err := h.authService.Register(req.Email, req.Username, req.Password)
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
