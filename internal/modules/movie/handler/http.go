package handler

import (
	"strconv"

	"github.com/geraldiaditya/ratix-backend/internal/modules/movie/service"
	"github.com/gofiber/fiber/v2"
)

type MovieHandler struct {
	Service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{Service: s}
}

func (h *MovieHandler) RegisterRoutes(app *fiber.App) {
	movies := app.Group("/movies")
	movies.Get("/categories", h.handleGetCategories)
	movies.Get("/banner", h.handleGetBanner)
	movies.Get("/", h.handleGetMovies) // List with query param
	movies.Get("/:id", h.handleDetail)
}

func (h *MovieHandler) handleGetCategories(c *fiber.Ctx) error {
	resp, err := h.Service.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(resp)
}

func (h *MovieHandler) handleGetBanner(c *fiber.Ctx) error {
	resp, err := h.Service.GetBanner()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if resp == nil {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.JSON(resp)
}

func (h *MovieHandler) handleGetMovies(c *fiber.Ctx) error {
	category := c.Query("category")
	resp, err := h.Service.GetMovies(category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(resp)
}

func (h *MovieHandler) handleDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	resp, err := h.Service.GetDetail(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(resp)
}
