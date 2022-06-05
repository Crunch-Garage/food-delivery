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

var restaurant_image = ""

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var restaurant models.Restaurant

	_ = json.NewDecoder(r.Body).Decode(&restaurant)

	var dbRestaurant models.Restaurant
	database.DB.Where("restaurant_name = ?", restaurant.Restaurant_name).First(&dbRestaurant)
	if dbRestaurant.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Business name already exists")
		return
	}

	restaurant.Registration_status = "PENDING"
	createdMenu := database.DB.Create(&restaurant)
	err = createdMenu.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMenu.Value)
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	var restaurants []models.Restaurant
	var restaurantsHolder []map[string]interface{}

	database.DB.Find(&restaurants)

	for i, _ := range restaurants {

		var profile models.Profile

		database.DB.Model(&restaurants[i]).Related(&profile)

		/*restaurant interface*/
		restaurantData := map[string]interface{}{
			"id":               restaurants[i].ID,
			"restaurant_image": restaurants[i].Restaurant_image,
			"restaurant_name":  restaurants[i].Restaurant_name,
			"phone_number":     restaurants[i].Phone_number,
			"address":          restaurants[i].Address,
			"location":         restaurants[i].LocationID,
			"owner":            profile,
		}

		restaurantsHolder = append(restaurantsHolder, restaurantData)

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurantsHolder)
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var restaurant models.Restaurant
	var profile models.Profile
	var location models.Location

	database.DB.First(&restaurant, id)
	database.DB.Model(&restaurant).Related(&profile)
	database.DB.Model(&restaurant).Related(&location)

	if restaurant.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Restaurant not found")
		return
	}

	restaurantData := map[string]interface{}{
		"id":               restaurant.ID,
		"CreatedAt":        restaurant.CreatedAt,
		"UpdatedAt":        restaurant.UpdatedAt,
		"restaurant_image": restaurant.Restaurant_image,
		"restaurant_name":  restaurant.Restaurant_name,
		"phone_number":     restaurant.Phone_number,
		"address":          restaurant.Address,
		"location":         location,
		"owner":            restaurant.ProfileID,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurantData)
}

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var restaurant models.Restaurant
	var dbRestaurant models.Restaurant
	var profile models.Profile
	var location models.Location

	database.DB.First(&dbRestaurant, id)
	database.DB.Model(&dbRestaurant).Related(&profile)
	database.DB.Model(&dbRestaurant).Related(&location)

	if dbRestaurant.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Restaurant not found")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&restaurant)

	file, _, _ := r.FormFile("restaurant_image")
	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "restaurant_image", config.EnvCloudMenuFolder())
		if err != nil {
			restaurant_image = dbRestaurant.Restaurant_image
		}
		restaurant_image = avatarUrl
	}

	database.DB.Model(&dbRestaurant).Updates(models.Restaurant{
		Restaurant_name:  restaurant.Restaurant_name,
		Address:          restaurant.Address,
		LocationID:       restaurant.LocationID,
		Phone_number:     restaurant.Phone_number,
		Restaurant_image: restaurant_image,
	})

	restaurantData := map[string]interface{}{
		"id":               dbRestaurant.ID,
		"restaurant_image": dbRestaurant.Restaurant_image,
		"restaurant_name":  dbRestaurant.Restaurant_name,
		"phone_number":     dbRestaurant.Phone_number,
		"address":          dbRestaurant.Address,
		"owner":            dbRestaurant.ProfileID,
		"location":         location,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurantData)
}
