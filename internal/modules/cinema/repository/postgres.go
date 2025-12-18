package repository

import (
	"github.com/geraldiaditya/ratix-backend/internal/modules/cinema/domain"
	"gorm.io/gorm"
)

type PostgresCinemaRepository struct {
	DB *gorm.DB
}

func NewPostgresCinemaRepository(db *gorm.DB) *PostgresCinemaRepository {
	return &PostgresCinemaRepository{DB: db}
}

func (r *PostgresCinemaRepository) GetAllCities() ([]string, error) {
	var cities []string
	if err := r.DB.Model(&domain.Cinema{}).Distinct("city").Pluck("city", &cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *PostgresCinemaRepository) GetCinemasByCity(city string) ([]domain.Cinema, error) {
	var cinemas []domain.Cinema
	if err := r.DB.Where("city = ?", city).Find(&cinemas).Error; err != nil {
		return nil, err
	}
	return cinemas, nil
}

func (r *PostgresCinemaRepository) GetByID(id int64) (*domain.Cinema, error) {
	var cinema domain.Cinema
	if err := r.DB.First(&cinema, id).Error; err != nil {
		return nil, err
	}
	return &cinema, nil
}

func (r *PostgresCinemaRepository) Create(cinema *domain.Cinema) error {
	return r.DB.Create(cinema).Error
}

func (r *PostgresCinemaRepository) GetCinemaByShowtimeID(showtimeID int64) (*domain.Cinema, error) {
	// Join Showtime and Cinema tables
	// Assuming tables are "showtimes" and "cinemas"
	// and showtimes has cinema_id
	var cinema domain.Cinema

	// We need to access Showtime table which is likely "showtimes".
	// Since we don't import movieDomain here, we can use raw query or map to struct if we had it.
	// But actually, we can just do a join if we know the schema.
	// "SELECT cinemas.* FROM cinemas JOIN showtimes ON showtimes.cinema_id = cinemas.id WHERE showtimes.id = ?"

	if err := r.DB.Joins("JOIN showtimes ON showtimes.cinema_id = cinemas.id").
		Where("showtimes.id = ?", showtimeID).
		First(&cinema).Error; err != nil {
		return nil, err
	}
	return &cinema, nil
}
