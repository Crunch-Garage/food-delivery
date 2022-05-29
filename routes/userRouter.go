package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router) {
	router.HandleFunc("/api/user/signup", controller.SignUp).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Login).Methods("POST")
	router.Handle("/api/user/{id}", middleware.IsAuthorized(controller.GetUser)).Methods("GET")
	router.Handle("/api/user/{id}", middleware.IsAuthorized(controller.UpdateUser)).Methods("PATCH")
}
