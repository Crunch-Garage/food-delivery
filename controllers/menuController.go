package controller

import (
	"crunchgarage/restaurant-food-delivery/config"
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var menu_image = ""

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	var menu models.Menu

	_ = json.NewDecoder(r.Body).Decode(&menu)

	createdMenu := database.DB.Create(&menu)
	err = createdMenu.Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMenu.Value)
}

func GetMenus(w http.ResponseWriter, r *http.Request) {

	var menus []models.Menu

	menuList := database.DB.Find(&menus)
	err = menuList.Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menus)
}

func GetMenu(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var menu models.Menu

	database.DB.First(&menu, id)

	if menu.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(menu)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menu)
}

func UpdateMenu(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var menu models.Menu
	var dbMenu models.Menu
	database.DB.First(&dbMenu, id)

	if dbMenu.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Menu not found")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&menu)

	if menu.Name != "" {
		dbMenu.Name = menu.Name
	}

	file, _, _ := r.FormFile("menu_image")
	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "menu_image", config.EnvCloudMenuFolder())
		if err != nil {
			menu_image = dbMenu.Menu_image
		}
		menu_image = avatarUrl
	}

	/*update menu*/
	updatedMenu := database.DB.Model(&menu).Updates(models.Menu{
		Name:       dbMenu.Name,
		Menu_image: menu_image,
	})
	err := updatedMenu.Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMenu.Value)

}
