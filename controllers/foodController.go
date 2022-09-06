package controller

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var food_image = ""

func CreateFood(w http.ResponseWriter, r *http.Request) {
	var food models.Food

	/*get formdata*/
	food_name := r.PostFormValue("name")
	food_price := r.PostFormValue("price")
	food_menu_id := r.PostFormValue("menu_id")
	food_restarant_id := r.PostFormValue("restarant_id")
	food_description := r.PostFormValue("description")

	file, _, _ := r.FormFile("food_image")

	if food_name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Food name is required")
		return
	}

	if food_price == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Food price is required")
		return
	}

	if food_menu_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Menu id is required")
		return
	}

	if food_restarant_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Restaurant id is required")
		return
	}

	if food_description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Food Description id is required")
		return
	}

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "food_image", config.EnvCloudFoodFolder())
		if err != nil {
			avatarUrl = ""
		}
		food_image = avatarUrl
	}

	food.Name = food_name
	food.Price, _ = strconv.ParseFloat(food_price, 64)
	food.Food_image = food_image
	food.MenuID, _ = strconv.Atoi(food_menu_id)
	food.RestaurantID, _ = strconv.Atoi(food_restarant_id)

	createdFood := database.DB.Create(&food)
	err = createdFood.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	food_image = ""
	w.WriteHeader(http.StatusCreated)
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Food not found")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&food)

	if food.Name != "" {
		dbFood.Name = food.Name
	}

	if strconv.FormatFloat(food.Price, 'E', -1, 32) != "" {
		dbFood.Price = food.Price
	}

	if strconv.Itoa(food.MenuID) != "" {
		dbFood.MenuID = food.MenuID
	}

	if strconv.Itoa(food.RestaurantID) != "" {
		dbFood.RestaurantID = food.RestaurantID
	}

	if food.Description != "" {
		dbFood.Description = food.Description
	}

	file, _, _ := r.FormFile("food_image")
	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "food_image", config.EnvCloudFoodFolder())
		if err != nil {
			avatarUrl = dbFood.Food_image
		}
		food_image = avatarUrl
	}

	/*update menu*/
	updatedFood := database.DB.Model(&dbFood).Updates(models.Food{
		Name:         dbFood.Name,
		Price:        dbFood.Price,
		Food_image:   food_image,
		MenuID:       dbFood.MenuID,
		RestaurantID: dbFood.RestaurantID,
		Description:  dbFood.Description,
	})
	err := updatedFood.Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	food_image = ""
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedFood.Value)

}
