package domain

import (
	"time"

	movieDomain "github.com/geraldiaditya/ratix-backend/internal/modules/movie/domain"
	userDomain "github.com/geraldiaditya/ratix-backend/internal/modules/user/domain"
)

type Ticket struct {
	ID          int64             `gorm:"primaryKey" json:"id"`
	UserID      int64             `gorm:"not null" json:"user_id"`
	User        userDomain.User   `gorm:"foreignKey:UserID" json:"-"`
	MovieID     int64             `gorm:"not null" json:"movie_id"`
	Movie       movieDomain.Movie `gorm:"foreignKey:MovieID" json:"movie"`
	ShowtimeID  int64             `gorm:"not null" json:"showtime_id"`
	BookingCode string            `gorm:"type:varchar(20);unique;not null" json:"booking_code"` // For QR
	Seats       string            `gorm:"type:varchar(50);not null" json:"seats"`               // e.g. "G14, G15"
	CinemaName  string            `gorm:"type:varchar(100);not null" json:"cinema_name"`        // e.g. "AMC Empire 25"
	TheaterName string            `gorm:"type:varchar(50);not null" json:"theater_name"`        // e.g. "Auditorium 12"
	Price       float64           `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      string            `gorm:"type:varchar(20);default:'active'" json:"status"` // active, history, cancelled
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type TicketRepository interface {
	GetByUserID(userID int64, status string) ([]Ticket, error)
	GetByID(id int64) (*Ticket, error)
	GetBookedSeats(showtimeID int64) ([]string, error)
	Create(ticket *Ticket) error
}
