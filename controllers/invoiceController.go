package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
	"time"
)

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice
	_ = json.NewDecoder(r.Body).Decode(&invoice)

	if invoice.OrderID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order id is required")
		return
	}

	if invoice.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Customer id is required")
		return
	}

	invoice_ := models.Invoice{
		UserID:         invoice.UserID,
		Amount:         invoice.Amount,
		OrderID:        invoice.OrderID,
		Payment_date:   time.Now().Format(time.RFC3339),
		Payment_method: "CARD",
		Payment_status: "PAID",
	}

	createdInvoice := database.DB.Create(&invoice_)
	err := createdInvoice.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInvoice.Value)
}
