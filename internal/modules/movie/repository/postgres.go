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
	if err := r.DB.Preload("Genres").Preload("Cast").Preload("Showtimes").First(&movie, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}

func (r *PostgresMovieRepository) GetByStatus(status string) ([]domain.Movie, error) {
	var movies []domain.Movie
	if err := r.DB.Where("status = ?", status).Preload("Genres").Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *PostgresMovieRepository) GetByGenre(genreName string) ([]domain.Movie, error) {
	var movies []domain.Movie
	// Join with movie_genres table and genres table to filter by genre name
	err := r.DB.Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Joins("JOIN genres ON genres.id = movie_genres.genre_id").
		Where("genres.name = ?", genreName).
		Preload("Genres").
		Find(&movies).Error
	if err != nil {
		return nil, err
	}
	return movies, nil
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
