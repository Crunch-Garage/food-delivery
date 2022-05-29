package models

import "github.com/jinzhu/gorm"

type Invoice struct {
	gorm.Model

	OrderID        int     `json:"order_id"`
	Payment_method string  `json:"payment_method" validate:"eq=CARD|eq=MOBILEMONEY"`
	Payment_status string  `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	Payment_date   string  `json:"payment_date"`
	UserID         int     `json:"customer_id"`
	Amount         float64 `json:"amount"`
}
