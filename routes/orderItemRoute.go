package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"

	"github.com/gorilla/mux"
)

func OrderItemRouter(router *mux.Router) {
	router.HandleFunc("/api/orderItems", controller.GetOrderItems).Methods("GET")
	router.HandleFunc("/api/orderItems/restaurant/{id}", controller.GetRestaurantOrderItems).Methods("GET")
}
