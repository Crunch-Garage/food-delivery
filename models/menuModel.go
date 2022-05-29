package models

import "github.com/jinzhu/gorm"

type Menu struct {
	gorm.Model

	Name       string `json:"name" validate:"required"`
	Menu_image string `json:"menu_image"`
}
