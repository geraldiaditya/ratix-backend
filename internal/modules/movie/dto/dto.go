package dto

import (
	"github.com/geraldiaditya/ratix-backend/internal/modules/movie/domain"
)

type GenreResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type BannerResponse struct {
	MovieID   int64    `json:"movie_id"`
	Title     string   `json:"title"`
	PosterURL string   `json:"poster_url"`
	Rating    float64  `json:"rating"`
	Genres    []string `json:"genres"`
}

type MovieListResponse struct {
	Movies []MovieResponse `json:"movies"`
}

type MovieResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"` // Optional for list view
	Duration    int      `json:"duration"`
	Rating      float64  `json:"rating"`
	PosterURL   string   `json:"poster_url"`
	Genres      []string `json:"genres,omitempty"`
}

type MovieDetailResponse struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Duration    int                `json:"duration"`
	Rating      float64            `json:"rating"`
	PosterURL   string             `json:"poster_url"`
	ReleaseDate string             `json:"release_date"`
	Genres      []string           `json:"genres"`
	Cast        []CastResponse     `json:"cast"`
	Showtimes   []ShowtimeResponse `json:"showtimes"`
}

type CastResponse struct {
	Name          string `json:"name"`
	Role          string `json:"role"`
	CharacterName string `json:"character_name"`
	PhotoURL      string `json:"photo_url"`
}

type ShowtimeResponse struct {
	StartTime string  `json:"start_time"`
	Price     float64 `json:"price"`
	Date      string  `json:"date"`
}

func ToMovieResponse(m domain.Movie) MovieResponse {
	genres := make([]string, len(m.Genres))
	for i, g := range m.Genres {
		genres[i] = g.Name
	}
	return MovieResponse{
		ID:        m.ID,
		Title:     m.Title,
		Duration:  m.Duration,
		Rating:    m.Rating,
		PosterURL: m.PosterURL,
		Genres:    genres,
	}
}

func ToMovieDetailResponse(m *domain.Movie) *MovieDetailResponse {
	genres := make([]string, len(m.Genres))
	for i, g := range m.Genres {
		genres[i] = g.Name
	}

	cast := make([]CastResponse, len(m.Cast))
	for i, c := range m.Cast {
		cast[i] = CastResponse{
			Name:          c.Name,
			Role:          c.Role,
			CharacterName: c.CharacterName,
			PhotoURL:      c.PhotoURL,
		}
	}

	showtimes := make([]ShowtimeResponse, len(m.Showtimes))
	for i, s := range m.Showtimes {
		showtimes[i] = ShowtimeResponse{
			StartTime: s.StartTime.Format("15:04"),
			Price:     s.Price,
			Date:      s.StartTime.Format("2006-01-02"),
		}
	}

	return &MovieDetailResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		Rating:      m.Rating,
		PosterURL:   m.PosterURL,
		ReleaseDate: m.ReleaseDate.Format("2006-01-02"),
		Genres:      genres,
		Cast:        cast,
		Showtimes:   showtimes,
	}
}
