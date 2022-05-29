package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gorilla/mux"
)

func OrderRouter(router *mux.Router) {
	router.Handle("/api/order/create", middleware.IsAuthorized(controller.CreateOrder)).Methods("POST")
	router.HandleFunc("/api/order/{id}", controller.GetOrder).Methods("GET")
	router.Handle("/api/order/{id}", middleware.IsAuthorized(controller.UpdateOrder)).Methods("PATCH")
}
