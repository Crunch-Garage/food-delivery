package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
)

func GetLocations(w http.ResponseWriter, r *http.Request) {

	var location []models.Location

	database.DB.Find(&location)

	json.NewEncoder(w).Encode(location)
}
