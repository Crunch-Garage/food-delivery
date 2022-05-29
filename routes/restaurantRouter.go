package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gorilla/mux"
)

func RestaurantRouter(router *mux.Router) {
	router.HandleFunc("/api/restaurants", controller.GetRestaurants).Methods("GET")
	router.HandleFunc("/api/restaurant/{id}", controller.GetRestaurant).Methods("GET")
	router.Handle("/api/restaurant/{id}", middleware.IsAuthorized(controller.UpdateRestaurant)).Methods("PATCH")
	router.Handle("/api/restaurants/create", middleware.IsAuthorized(controller.CreateRestaurant)).Methods("POST")

}
