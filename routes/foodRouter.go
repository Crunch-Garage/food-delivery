package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gorilla/mux"
)

func FoodRouter(router *mux.Router) {
	router.HandleFunc("/api/foods", controller.GetFoods).Methods("GET")
	router.HandleFunc("/api/food/{id}", controller.GetFood).Methods("GET")
	router.Handle("/api/food/{id}", middleware.IsAuthorized(controller.UpdateFood)).Methods("PATCH")
	router.Handle("/api/foods", middleware.IsAuthorized(controller.CreateFood)).Methods("POST")

}
