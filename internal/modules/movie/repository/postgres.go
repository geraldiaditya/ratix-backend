package repository

import (
	"errors"
	"fmt"

	"github.com/geraldiaditya/ratix-backend/internal/modules/movie/domain"
	"gorm.io/gorm"
)

type PostgresMovieRepository struct {
	DB *gorm.DB
}

func NewPostgresMovieRepository(db *gorm.DB) *PostgresMovieRepository {
	return &PostgresMovieRepository{DB: db}
}

func (r *PostgresMovieRepository) GetAll() ([]domain.Movie, error) {
	var movies []domain.Movie
	if err := r.DB.Preload("Genres").Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *PostgresMovieRepository) GetByID(id int64) (*domain.Movie, error) {
	var movie domain.Movie
	if err := r.DB.Preload("Genres").Preload("Cast").Preload("Showtimes.Cinema").First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}

func (r *PostgresMovieRepository) GetByStatus(status string, limit, offset int) ([]domain.Movie, int64, error) {
	var movies []domain.Movie
	var total int64

	// Count total
	if err := r.DB.Model(&domain.Movie{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.DB.Where("status = ?", status).Limit(limit).Offset(offset).Preload("Genres").Find(&movies).Error; err != nil {
		return nil, 0, err
	}
	return movies, total, nil
}

func (r *PostgresMovieRepository) GetByGenre(genreName string, limit, offset int) ([]domain.Movie, int64, error) {
	var movies []domain.Movie
	var total int64

	// Count total (need proper join for count too)
	err := r.DB.Model(&domain.Movie{}).
		Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Joins("JOIN genres ON genres.id = movie_genres.genre_id").
		Where("genres.name = ?", genreName).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Fetch data
	err = r.DB.Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Joins("JOIN genres ON genres.id = movie_genres.genre_id").
		Where("genres.name = ?", genreName).
		Limit(limit).Offset(offset).
		Preload("Genres").
		Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}
	return movies, total, nil
}

func (r *PostgresMovieRepository) GetAllGenres() ([]domain.Genre, error) {
	var genres []domain.Genre
	if err := r.DB.Find(&genres).Error; err != nil {
		return nil, err
	}
	return genres, nil
}

func (r *PostgresMovieRepository) Create(movie *domain.Movie) error {
	return r.DB.Create(movie).Error
}
