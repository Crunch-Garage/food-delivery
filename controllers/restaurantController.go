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

	createdMenu := database.DB.Create(&restaurant)
	err = createdMenu.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(createdMenu.Value)
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	var restaurants []models.Restaurant
	var restaurantsHolder []map[string]interface{}

	database.DB.Find(&restaurants)

	for i, _ := range restaurants {

		var user models.User

		database.DB.Model(&restaurants[i]).Related(&user)

		/*user interface*/
		userData := map[string]interface{}{
			"id":           user.ID,
			"first_name":   user.First_name,
			"last_name":    user.Last_name,
			"user_name":    user.User_name,
			"email":        user.Email,
			"avatar":       user.Avatar,
			"phone":        user.Phone,
			"account_type": user.Account_type,
		}

		/*restaurant interface*/
		restaurantData := map[string]interface{}{
			"id":               restaurants[i].ID,
			"restaurant_image": restaurants[i].Restaurant_image,
			"restaurant_name":  restaurants[i].Restaurant_name,
			"phone_number":     restaurants[i].Phone_number,
			"address":          restaurants[i].Address,
			"location":         restaurants[i].Location,
			"owner":            userData,
		}

		restaurantsHolder = append(restaurantsHolder, restaurantData)

	}

	json.NewEncoder(w).Encode(restaurantsHolder)
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var restaurant models.Restaurant
	var user models.User

	database.DB.First(&restaurant, id)
	database.DB.Model(&restaurant).Related(&user)

	if restaurant.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Restaurant not found")
		return
	}

	/*user interface*/
	userData := map[string]interface{}{
		"id":           user.ID,
		"first_name":   user.First_name,
		"last_name":    user.Last_name,
		"user_name":    user.User_name,
		"email":        user.Email,
		"avatar":       user.Avatar,
		"phone":        user.Phone,
		"account_type": user.Account_type,
	}

	restaurantData := map[string]interface{}{
		"id":               restaurant.ID,
		"CreatedAt":        restaurant.CreatedAt,
		"UpdatedAt":        restaurant.UpdatedAt,
		"restaurant_image": restaurant.Restaurant_image,
		"restaurant_name":  restaurant.Restaurant_name,
		"phone_number":     restaurant.Phone_number,
		"address":          restaurant.Address,
		"location":         restaurant.Location,
		"owner":            userData,
	}

	/*restaurant interface*/
	json.NewEncoder(w).Encode(restaurantData)
}

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var restaurant models.Restaurant
	var dbRestaurant models.Restaurant
	database.DB.First(&dbRestaurant, id)

	if dbRestaurant.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Restaurant not found")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&restaurant)

	file, _, _ := r.FormFile("restaurant_image")
	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "restaurant_image")
		if err != nil {
			dbRestaurant.Restaurant_image = restaurant.Restaurant_image
		}
		dbRestaurant.Restaurant_image = avatarUrl
		dbRestaurant.Restaurant_name = restaurant.Restaurant_name

	}

	dbRestaurant.Restaurant_name = restaurant.Restaurant_name
	dbRestaurant.Address = restaurant.Address
	dbRestaurant.Location = restaurant.Location
	dbRestaurant.Phone_number = restaurant.Phone_number

	database.DB.Save(&dbRestaurant)

	json.NewEncoder(w).Encode(dbRestaurant)
}
