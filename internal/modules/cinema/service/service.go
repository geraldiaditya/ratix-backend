package service

import (
	"fmt"
	"strings"

	"github.com/geraldiaditya/ratix-backend/internal/modules/cinema/domain"
	"github.com/geraldiaditya/ratix-backend/internal/modules/cinema/dto"
	ticketDomain "github.com/geraldiaditya/ratix-backend/internal/modules/ticket/domain"
)

type CinemaService struct {
	Repo       domain.CinemaRepository
	TicketRepo ticketDomain.TicketRepository
}

func NewCinemaService(repo domain.CinemaRepository, ticketRepo ticketDomain.TicketRepository) *CinemaService {
	return &CinemaService{Repo: repo, TicketRepo: ticketRepo}
}

func (s *CinemaService) GetLocations() (*dto.CityResponse, error) {
	cities, err := s.Repo.GetAllCities()
	if err != nil {
		return nil, err
	}
	return &dto.CityResponse{Cities: cities}, nil
}

func (s *CinemaService) GetCinemas(city string) ([]dto.CinemaResponse, error) {
	cinemas, err := s.Repo.GetCinemasByCity(city)
	if err != nil {
		return nil, err
	}

	var resp []dto.CinemaResponse
	for _, c := range cinemas {
		resp = append(resp, dto.ToCinemaResponse(c))
	}
	return resp, nil
}

func (s *CinemaService) GetSeatLayout(showtimeID int64) (*dto.SeatLayoutResponse, error) {
	// 1. Fetch Booked Seats
	bookedSeatStrings, err := s.TicketRepo.GetBookedSeats(showtimeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch booked seats: %w", err)
	}

	// 1b. Fetch Cinema details via Showtime (we need to know which cinema to get base price)
	// Currently we only have showtimeID. We need to fetch Showtime -> Cinema.
	// We don't have ShowtimeRepo here. But we can fetch Cinema if we had CinemaID.
	// Ideally TicketRepo or a new ShowtimeRepo method would give us this.
	// For now, let's assume we can query Cinema directly if we knew the ID.
	// But we don't.
	// HACK: We should really inject ShowtimeRepo or similar.
	// OR: We can just fetch the Cinema associated with the showtime using GORM if we had access to DB.
	// s.Repo is CinemaRepository. It can get Cinema by ID. But we need CinemaID from ShowtimeID.

	// Let's rely on TicketRepo to give us the cinema ID? Or just fetch everything?
	// To do this properly, we need to fetch the Showtime entity first.
	// Since we don't have ShowtimeRepo injected, we can't easily do it without changing dependencies.

	// However, the prompt says "Refactor...".
	// I should probably inject ShowtimeRepo or similar?
	// OR: I can modify TicketRepo to return the Cinema ID for the showtime? NO, that's weird.

	// Better approach:
	// Use `DB` directly? No, abstraction.

	// Let's check `CinemaRepository`. Maybe add `GetCinemaByShowtimeID`?
	// `GetCinemasByCity`, `GetByID`.

	// I'll add `GetCinemaByShowtimeID` to `CinemaRepository` interface? Easiest way.
	// But CinemaRepo is for Cinema domain.
	// A showtime belongs to a cinema. So `GetCinemaByShowtimeID` is valid in CinemaRepo?
	// It crosses domains slightly but it's okay.

	// Or, I can leave the price hardcoded for a sec while I fix the lints? No, user wants logic.

	// Let's assume for this specific task I can add `GetCinemaByShowtimeID(showtimeID int64)` to `CinemaRepository`.
	cinema, err := s.Repo.GetCinemaByShowtimeID(showtimeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cinema for showtime: %w", err)
	}

	// flatten booked seats: "A1, A2" -> ["A1", "A2"]
	bookedMap := make(map[string]bool)
	for _, s := range bookedSeatStrings {
		parts := strings.Split(s, ",")
		for _, p := range parts {
			bookedMap[strings.TrimSpace(p)] = true
		}
	}

	// 2. Generate Static Layout (Row A-L, Cols 1-8 for example)
	// In real app, this should come from Theater layout in DB.
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	cols := 8
	var seats []dto.Seat

	for i, r := range rows {
		seatType := "standard"
		price := cinema.BasePrice
		// Rows A-D are standard? User image shows Rows A-L.
		// Let's assume K, L are premium or maybe closest to screen?
		// Actually typical cinema: Back is VIP/Premium.
		// Let's make J, K, L premium.
		if i >= 9 { // J, K, L
			seatType = "premium"
			price += 25000.0 // Premium surcharge
		}

		for c := 1; c <= cols; c++ {
			seatNum := fmt.Sprintf("%s%d", r, c)
			status := "available"
			if bookedMap[seatNum] {
				status = "occupied"
			}

			seats = append(seats, dto.Seat{
				Row:    r,
				Number: c,
				Status: status,
				Type:   seatType,
				Price:  price,
			})
		}
	}

	return &dto.SeatLayoutResponse{
		Layout: dto.SeatLayout{
			Rows:  len(rows),
			Cols:  cols,
			Seats: seats,
		},
		Legend: dto.SeatLegend{
			Available: "Available",
			Occupied:  "Occupied",
			Selected:  "Selected",
		},
	}, nil
}
