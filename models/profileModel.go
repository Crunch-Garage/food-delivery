package models

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model

	UserID        int          `json:"user_id"`
	First_name    string       `json:"first_name"`
	Last_name     string       `json:"last_name"`
	Profile_image string       `json:"profile_image"`
	User_type     string       `json:"user_type" validate:"eq=PRO|eq=CLIENT"`
	Pro_type      string       `json:"pro_type"  validate:"eq=CHEF|eq=RIDER"`
	User_name     string       `json:"user_name"`
	Restaurant    []Restaurant `json:"restaurant"`
}
