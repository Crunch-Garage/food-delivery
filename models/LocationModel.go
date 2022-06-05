package models

import "github.com/jinzhu/gorm"

type Location struct {
	gorm.Model

	Name      string  `json:"location_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var (
	Locations = []Location{
		{Name: "Kilifi", Latitude: -3.6225519, Longitude: 39.8474175},
		{Name: "Nairobi", Latitude: -1.3031933, Longitude: 36.5664722},
		{Name: "Mombasa", Latitude: -4.0351767, Longitude: 39.5258227},
	}
)
