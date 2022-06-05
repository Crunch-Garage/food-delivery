package database

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/models"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

var err error
var DB *gorm.DB

func OpenDB() *gorm.DB {

	/*
		loading environmental variables
	*/
	dialect := config.EnvDBDialect()
	host := config.EnvDBHost()
	dbport := config.EnvDBPort()
	user := config.EnvDBUser()
	dbName := config.EnvDBName()
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
		&models.Profile{},
		&models.Menu{},
		&models.Restaurant{},
		&models.Food{},
		&models.Order{},
		&models.OrderItem{},
		&models.Invoice{},
		&models.Location{},
	)
}
