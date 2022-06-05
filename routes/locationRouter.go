package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"

	"github.com/gorilla/mux"
)

func LocationRouter(router *mux.Router) {
	router.HandleFunc("/api/location", controller.GetLocations).Methods("GET")

}
