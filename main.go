package main

//jwt "github.com/dgrijalva/jwt-go"
import (
	"log"
	"net/http"

	"crunchgarage/restaurant-food-delivery/database"
	routes "crunchgarage/restaurant-food-delivery/routes"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	database.OpenDB()

	defer database.CloseDB()

	database.AutoMigrate()

	handleRequests()
}

/* Handle API requests*/
func handleRequests() {

	router := mux.NewRouter()

	routes.UserRouter(router)
	routes.MenuRouter(router)
	routes.RestaurantRouter(router)
	routes.FoodRouter(router)
	routes.OrderRouter(router)
	routes.OrderItemRouter(router)
	routes.InvoiceRouter(router)
	routes.LocationRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
