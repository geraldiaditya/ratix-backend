package dto

import "github.com/geraldiaditya/ratix-backend/internal/modules/cinema/domain"

type CityResponse struct {
	Cities []string `json:"cities"`
}

type CinemaResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type SeatLayoutResponse struct {
	Layout SeatLayout `json:"layout"`
	Legend SeatLegend `json:"legend"`
}

type SeatLayout struct {
	Rows  int    `json:"rows"`
	Cols  int    `json:"cols"`
	Seats []Seat `json:"seats"`
}

type Seat struct {
	Row    string  `json:"row"`
	Number int     `json:"number"`
	Status string  `json:"status"` // available, occupied, selected
	Type   string  `json:"type"`   // standard, premium
	Price  float64 `json:"price"`
}

type SeatLegend struct {
	Available string `json:"available"`
	Occupied  string `json:"occupied"`
	Selected  string `json:"selected"`
}

func ToCinemaResponse(c domain.Cinema) CinemaResponse {
	return CinemaResponse{
		ID:      c.ID,
		Name:    c.Name,
		City:    c.City,
		Address: c.Address,
	}
}
