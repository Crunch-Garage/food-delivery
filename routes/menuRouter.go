package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gorilla/mux"
)

func MenuRouter(router *mux.Router) {
	router.HandleFunc("/api/menus", controller.GetMenus).Methods("GET")
	router.HandleFunc("/api/menu/{id}", controller.GetMenu).Methods("GET")
	router.Handle("/api/menu/{id}", middleware.IsAuthorized(controller.UpdateMenu)).Methods("PATCH")
	router.Handle("/api/menus", middleware.IsAuthorized(controller.CreateMenu)).Methods("POST")
}
