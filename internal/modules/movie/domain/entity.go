package domain

import (
	"time"
)

type Movie struct {
	ID          int64        `gorm:"primaryKey" json:"id"`
	Title       string       `gorm:"not null;type:varchar(255)" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Duration    int          `gorm:"not null" json:"duration"` // in minutes
	Rating      float64      `gorm:"type:decimal(3,1)" json:"rating"`
	PosterURL   string       `gorm:"type:varchar(255)" json:"poster_url"`
	ReleaseDate time.Time    `gorm:"type:date" json:"release_date"`
	Status      string       `gorm:"type:varchar(50);default:'now_showing'" json:"status"` // now_showing, coming_soon
	Genres      []Genre      `gorm:"many2many:movie_genres;" json:"genres"`
	Cast        []CastMember `gorm:"foreignKey:MovieID" json:"cast"`
	Showtimes   []Showtime   `gorm:"foreignKey:MovieID" json:"showtimes"`
}

type Genre struct {
	ID   int64  `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null;unique;type:varchar(100)" json:"name"`
}

type CastMember struct {
	ID            int64  `gorm:"primaryKey" json:"id"`
	MovieID       int64  `gorm:"not null" json:"movie_id"`
	Name          string `gorm:"not null;type:varchar(255)" json:"name"`
	Role          string `gorm:"not null;type:varchar(50)" json:"role"`   // Actor, Director, etc.
	CharacterName string `gorm:"type:varchar(255)" json:"character_name"` // For actors
	PhotoURL      string `gorm:"type:varchar(255)" json:"photo_url"`
}

type Showtime struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	MovieID   int64     `gorm:"not null" json:"movie_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	Price     float64   `gorm:"not null;type:decimal(10,2)" json:"price"`
}

type MovieRepository interface {
	GetAll() ([]Movie, error)
	GetByID(id int64) (*Movie, error)
	GetByStatus(status string) ([]Movie, error)
	GetByGenre(genre string) ([]Movie, error)
	GetAllGenres() ([]Genre, error)
	Create(movie *Movie) error
}
