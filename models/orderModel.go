package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model

	UserID           int         `json:"customer_id"`
	OrderItem        []OrderItem `json:"order_items"`
	Delivery_address string      `json:"delivery_address"`
	Order_status     string      `json:"order_status" validate:"eq=PENDING|eq=CANCELLED|eq=DELIVERED"`
	DriverID         int         `json:"driver_id"`
	Order_Date       string      `json:"order_date"`
	Total_price      float64     `json:"total_price"`     //total coast of items
	Delivery_charge  float64     `json:"delivery_charge"` // estimated delivery charge
	Total_amount     float64     `json:"total_amount"`    // total amount to be paid for order and delivery charge
}
