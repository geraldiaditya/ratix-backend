package dto

import (
	"time"

	"github.com/geraldiaditya/ratix-backend/internal/modules/ticket/domain"
)

type TicketListResponse struct {
	Tickets []TicketResponse `json:"tickets"`
}

type TicketResponse struct {
	ID         int64     `json:"id"`
	MovieTitle string    `json:"movie_title"`
	PosterURL  string    `json:"poster_url"`
	Date       time.Time `json:"date"` // Extracted from some logic or just CreatedAt? Ideally from Showtime
	Time       string    `json:"time"` // Formatted time
	CinemaName string    `json:"cinema_name"`
	IsActive   bool      `json:"is_active"`
}

type TicketDetailResponse struct {
	ID             int64   `json:"id"`
	MovieTitle     string  `json:"movie_title"`
	PosterURL      string  `json:"poster_url"`
	Rating         string  `json:"rating"`    // e.g. "PG-13 | 2h 46m"
	Score          string  `json:"score"`     // "93% Rotten Tomatoes"
	DateTimeString string  `json:"date_time"` // "Saturday, November 16, 2024 at 7:30 PM"
	CinemaName     string  `json:"cinema_name"`
	TheaterName    string  `json:"theater_name"`
	Seats          string  `json:"seats"`
	BookingCode    string  `json:"booking_code"` // QR content
	Price          float64 `json:"price"`
}

func ToTicketResponse(t domain.Ticket) TicketResponse {
	// Dummy date/time logic if Showtime relation isn't full populated
	// Ideally we use t.Showtime.StartTime

	return TicketResponse{
		ID:         t.ID,
		MovieTitle: t.Movie.Title,
		PosterURL:  t.Movie.PosterURL,
		Date:       t.CreatedAt, // Fallback
		Time:       t.CreatedAt.Format("15:04"),
		CinemaName: t.CinemaName,
		IsActive:   t.Status == "active",
	}
}

func ToTicketDetailResponse(t domain.Ticket) TicketDetailResponse {
	return TicketDetailResponse{
		ID:             t.ID,
		MovieTitle:     t.Movie.Title,
		PosterURL:      t.Movie.PosterURL,
		Rating:         "PG-13 | 2h 30m",                                     // Dummy
		Score:          "95% Rotten Tomatoes",                                // Dummy
		DateTimeString: t.CreatedAt.Format("Monday, 02 Jan 2006 at 3:04 PM"), // Placeholder
		CinemaName:     t.CinemaName,
		TheaterName:    t.TheaterName,
		Seats:          t.Seats,
		BookingCode:    t.BookingCode,
		Price:          t.Price,
	}
}
