package models

import "github.com/jinzhu/gorm"

type Restaurant struct {
	gorm.Model

	Restaurant_image    string `json:"restaurant_image"`
	Restaurant_name     string `json:"restaurant_name"`
	Phone_number        string `json:"phone_number"`
	Address             string `json:"address"`
	Location            string `json:"location"` // geo coordnates, change this interface
	ProfileID           int    `json:"owner"`
	Registration_status string `json:"registration_status" validate:"eq=PENDING|eq=ACCEPTED|eq=REJECTED"`
}
