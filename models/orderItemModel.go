package models

import "github.com/jinzhu/gorm"

type OrderItem struct {
	gorm.Model

	Quantity     int     `json:"quantity"`
	Unit_price   float64 `json:"unit_price"`
	OrderID      int     `json:"order_id"`
	FoodID       int     `json:"food_id"`
	RestaurantID int     `json:"restaurant_id"`
}
