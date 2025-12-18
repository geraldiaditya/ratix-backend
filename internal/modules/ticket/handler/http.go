package handler

import (
	"strconv"

	"github.com/geraldiaditya/ratix-backend/internal/modules/ticket/service"
	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	Service *service.TicketService
}

func NewTicketHandler(service *service.TicketService) *TicketHandler {
	return &TicketHandler{Service: service}
}

func (h *TicketHandler) RegisterRoutes(app *fiber.App) {
	tickets := app.Group("/tickets")
	// TODO: Add Middleware to get UserID
	tickets.Get("/", h.handleGetMyTickets)
	tickets.Get("/:id", h.handleGetTicketDetail)
}

func (h *TicketHandler) handleGetMyTickets(c *fiber.Ctx) error {
	// Mock UserID for now (should come from context/JWT)
	userID := int64(1)

	status := c.Query("status")
	resp, err := h.Service.GetMyTickets(userID, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(resp)
}

func (h *TicketHandler) handleGetTicketDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	resp, err := h.Service.GetTicketDetail(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return c.JSON(resp)
}
