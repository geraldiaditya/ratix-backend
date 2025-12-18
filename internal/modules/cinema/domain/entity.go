package domain

type Cinema struct {
	ID        int64   `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"not null;type:varchar(100)" json:"name"`
	City      string  `gorm:"not null;type:varchar(50)" json:"city"`
	Address   string  `gorm:"type:text" json:"address"`
	BasePrice float64 `gorm:"not null;type:decimal(10,2);default:50000" json:"base_price"`
}

type Theater struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	CinemaID int64  `gorm:"not null" json:"cinema_id"`
	Cinema   Cinema `gorm:"foreignKey:CinemaID" json:"cinema"`
	Name     string `gorm:"not null;type:varchar(50)" json:"name"` // e.g. "Studio 1", "IMAX"
	Type     string `gorm:"type:varchar(20)" json:"type"`          // Regular, IMAX, Premiere
}

type CinemaRepository interface {
	GetAllCities() ([]string, error)
	GetCinemasByCity(city string) ([]Cinema, error)
	GetByID(id int64) (*Cinema, error)
	GetCinemaByShowtimeID(showtimeID int64) (*Cinema, error)
	Create(cinema *Cinema) error
}
