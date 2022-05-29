package database

import (
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var err error
var DB *gorm.DB

func OpenDB() *gorm.DB {

	/*
		loading environmental variables
	*/
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbport := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := "zoom20$$" //os.Getenv("PASSWD")

	/*
		Database connection string
	*/
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, dbport, user, dbName, password)

	/*
		opening database connection
	*/
	DB, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)

	} else {
		fmt.Println("Successfully connected to database")
	}

	return DB

}

func CloseDB() error {
	return DB.Close()
}

func AutoMigrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Menu{},
		&models.Restaurant{},
		&models.Food{},
		&models.Order{},
		&models.OrderItem{},
		&models.Invoice{})
}
