package handler

import (
	"strconv"

	"github.com/geraldiaditya/ratix-backend/internal/modules/user/dto"
	"github.com/geraldiaditya/ratix-backend/internal/modules/user/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Service   *service.UserService
	Validator *validator.Validate
}

func NewUserHandler(s *service.UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{Service: s, Validator: v}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", h.handleRegister)
	auth.Post("/login", h.handleLogin)

	users := app.Group("/users")
	users.Get("/get", h.handleGetUser)
}

func (h *UserHandler) handleLogin(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.Validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// User Secret is now injected in Service
	resp, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	return c.JSON(resp)
}

func (h *UserHandler) handleRegister(c *fiber.Ctx) error {
	var req dto.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.Validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, err := h.Service.RegisterUser(req.Name, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(dto.ToUserResponse(user))
}

func (h *UserHandler) handleGetUser(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	user, err := h.Service.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(dto.ToUserResponse(user))
}
