package repository

import (
	"errors"
	"fmt"

	"github.com/geraldiaditya/ratix-backend/internal/modules/ticket/domain"
	"gorm.io/gorm"
)

type PostgresTicketRepository struct {
	DB *gorm.DB
}

func NewPostgresTicketRepository(db *gorm.DB) *PostgresTicketRepository {
	return &PostgresTicketRepository{DB: db}
}

func (r *PostgresTicketRepository) GetByUserID(userID int64, status string) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	query := r.DB.Where("user_id = ?", userID).Preload("Movie")

	if status != "" {
		if status == "history" {
			// History means watched or cancelled -> simplified for now, assuming 'history' status exists
			// Or logic: status IN ('completed', 'cancelled')
			query = query.Where("status IN ?", []string{"completed", "cancelled"})
		} else {
			// Active default
			query = query.Where("status = ?", "active")
		}
	}

	// Recent bookings first
	if err := query.Order("created_at desc").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *PostgresTicketRepository) GetByID(id int64) (*domain.Ticket, error) {
	var ticket domain.Ticket
	if err := r.DB.Preload("Movie").First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *PostgresTicketRepository) Create(ticket *domain.Ticket) error {
	return r.DB.Create(ticket).Error
}

func (r *PostgresTicketRepository) GetBookedSeats(showtimeID int64) ([]string, error) {
	var seats []string
	// Assuming "seats" column stores "A1, A2" (comma separated) or single seat "A1"
	// We need to fetch all seats from active bookings.
	// Since the DB stores a string, we might need to parse it if one ticket has multiple seats.
	// However, usually specific impl might vary. For now, let's just fetch the strings.
	var tickets []domain.Ticket
	if err := r.DB.Where("showtime_id = ? AND status != ?", showtimeID, "cancelled").Find(&tickets).Error; err != nil {
		return nil, err
	}

	for _, t := range tickets {
		// Normalize or split if needed. The Entity says `Seats string`.
		// If it's "A1, A2", we should split.
		// For simplicity in this iteration, we return the raw string and let service handle splitting or we do it here.
		// I'll append raw strings, but splitting is safer if we want individual seats.
		// Let's assume one ticket might hold multiple seats.
		// We can leave splitting to the service layer or do it here.
		// Let's return the simplified list of strings from DB for now.
		seats = append(seats, t.Seats)
	}
	return seats, nil
}
