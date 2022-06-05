package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	User_name     string    `json:"user_name"`
	Password      string    `json:"password"`
	Email         string    `gorm:"type:varchar(100);unique_index" json:"email"`
	Phone         string    `gorm:"type:varchar(100);unique_index" json:"phone"`
	First_name    string    `json:"first_name"`
	Last_name     string    `json:"last_name"`
	User_type     string    `json:"user_type" validate:"eq=PRO|eq=CLIENT"`
	Profile_image string    `json:"profile_image"`
	Profile       []Profile `gorm:"profile"`
}
