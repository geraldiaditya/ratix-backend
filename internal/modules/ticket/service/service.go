package service

import (
	"github.com/geraldiaditya/ratix-backend/internal/modules/ticket/domain"
	"github.com/geraldiaditya/ratix-backend/internal/modules/ticket/dto"
)

type TicketService struct {
	Repo domain.TicketRepository
}

func NewTicketService(repo domain.TicketRepository) *TicketService {
	return &TicketService{Repo: repo}
}

func (s *TicketService) GetMyTickets(userID int64, status string) (*dto.TicketListResponse, error) {
	tickets, err := s.Repo.GetByUserID(userID, status)
	if err != nil {
		return nil, err
	}

	var ticketResps []dto.TicketResponse
	for _, t := range tickets {
		ticketResps = append(ticketResps, dto.ToTicketResponse(t))
	}

	return &dto.TicketListResponse{Tickets: ticketResps}, nil
}

func (s *TicketService) GetTicketDetail(id int64) (*dto.TicketDetailResponse, error) {
	ticket, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToTicketDetailResponse(*ticket)
	return &resp, nil
}
