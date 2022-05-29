package models

import "github.com/jinzhu/gorm"

type Food struct {
	gorm.Model

	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Food_image   string  `json:"food_image"`
	MenuID       int     `json:"menu_id"`
	RestaurantID int     `json:"restarant_id"`
	Status       bool    `json:"status"`
}
