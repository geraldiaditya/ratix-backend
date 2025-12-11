package service

import (
	"github.com/geraldiaditya/ratix-backend/internal/modules/movie/domain"
	"github.com/geraldiaditya/ratix-backend/internal/modules/movie/dto"
)

type MovieService struct {
	Repo domain.MovieRepository
}

func NewMovieService(repo domain.MovieRepository) *MovieService {
	return &MovieService{Repo: repo}
}

func (s *MovieService) GetCategories() ([]dto.GenreResponse, error) {
	genres, err := s.Repo.GetAllGenres()
	if err != nil {
		return nil, err
	}
	var resp []dto.GenreResponse
	for _, g := range genres {
		resp = append(resp, dto.GenreResponse{ID: g.ID, Name: g.Name})
	}
	return resp, nil
}

func (s *MovieService) GetBanner() (*dto.BannerResponse, error) {
	// Logic: Get 'now_showing' AND standard picking logic (e.g. highest rated or first)
	movies, err := s.Repo.GetByStatus("now_showing")
	if err != nil {
		return nil, err
	}
	if len(movies) == 0 {
		return nil, nil // Or default banner
	}
	m := movies[0]

	genres := make([]string, len(m.Genres))
	for i, g := range m.Genres {
		genres[i] = g.Name
	}

	return &dto.BannerResponse{
		MovieID:   m.ID,
		Title:     m.Title,
		PosterURL: m.PosterURL,
		Rating:    m.Rating,
		Genres:    genres,
	}, nil
}

func (s *MovieService) GetMovies(category string) (*dto.MovieListResponse, error) {
	var movies []domain.Movie
	var err error

	if category == "" || category == "now_showing" || category == "coming_soon" {
		status := category
		if status == "" {
			status = "now_showing" // Default
		}
		movies, err = s.Repo.GetByStatus(status)
	} else {
		// Assume it's a genre
		movies, err = s.Repo.GetByGenre(category)
	}

	if err != nil {
		return nil, err
	}

	resp := make([]dto.MovieResponse, len(movies))
	for i, m := range movies {
		resp[i] = dto.ToMovieResponse(m)
	}

	return &dto.MovieListResponse{Movies: resp}, nil
}

func (s *MovieService) GetDetail(id int64) (*dto.MovieDetailResponse, error) {
	movie, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return dto.ToMovieDetailResponse(movie), nil
}
