package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateFood(w http.ResponseWriter, r *http.Request) {
	var food models.Food

	_ = json.NewDecoder(r.Body).Decode(&food)

	if food.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Name field is required")
		return
	}

	createdFood := database.DB.Create(&food)
	err = createdFood.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(createdFood.Value)
}

func GetFoods(w http.ResponseWriter, r *http.Request) {
	var foods []models.Food

	database.DB.Find(&foods)

	json.NewEncoder(w).Encode(foods)
}

func GetFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var food models.Food

	database.DB.First(&food, id)

	if food.ID == 0 {
		json.NewEncoder(w).Encode(food)
	}

	json.NewEncoder(w).Encode(food)
}

func UpdateFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var food models.Food
	var dbFood models.Food

	database.DB.First(&dbFood, id)

	if dbFood.ID == 0 {
		json.NewEncoder(w).Encode(dbFood)
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&food)

	file, _, _ := r.FormFile("food_image")
	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "food_image")
		if err != nil {
			dbFood.Food_image = food.Food_image
		}
		dbFood.Food_image = avatarUrl
	}

	dbFood.Name = food.Name
	dbFood.Description = food.Description
	dbFood.MenuID = food.MenuID
	dbFood.Price = food.Price

	database.DB.Save(&dbFood)

	json.NewEncoder(w).Encode(dbFood)

}
