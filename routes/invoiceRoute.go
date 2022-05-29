package routes

import (
	controller "crunchgarage/restaurant-food-delivery/controllers"
	"crunchgarage/restaurant-food-delivery/middleware"

	"github.com/gorilla/mux"
)

func InvoiceRouter(router *mux.Router) {
	router.Handle("/api/invoice/create", middleware.IsAuthorized(controller.CreateInvoice)).Methods("POST")
}
