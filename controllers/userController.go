package controller

import (
	"crunchgarage/restaurant-food-delivery/database"
	helper "crunchgarage/restaurant-food-delivery/helpers"
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var err error

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Password does not match")
		check = false
	}

	return check, msg
}

func SignUp(w http.ResponseWriter, r *http.Request) {

	var user models.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	var dbUser models.User
	database.DB.Where("email = ?", user.Email).First(&dbUser)
	if dbUser.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Email address already exists")
		return
	}

	database.DB.Where("phone = ?", user.Phone).First(&dbUser)
	if dbUser.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Phone number already exists")
		return
	}

	user_ := models.User{
		First_name:   user.First_name,
		Last_name:    user.Last_name,
		Password:     HashPassword(user.Password),
		Email:        user.Email,
		Avatar:       user.Avatar,
		Phone:        user.Phone,
		User_name:    strings.Split(user.Email, "@")[0], //split the email address at the @
		Account_type: "PRO",                             //CLIENT or PRO
	}

	createdUser := database.DB.Create(&user_)
	err = createdUser.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(createdUser)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User
	var dbUser models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	//find a user with username and see if that user even exists
	database.DB.Where("user_name = ?", user.User_name).First(&dbUser)

	if dbUser.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("User does not exist")
		return
	}

	//check if the password is correct
	passwordIsCorrect, msg := VerifyPassword(user.Password, dbUser.Password)
	if !passwordIsCorrect {
		json.NewEncoder(w).Encode(msg)
		return
	}

	accessTokenExpiration := time.Duration(60) * time.Minute
	refreshTokenExpiration := time.Duration(30*24) * time.Hour

	accessToken, accessTokenExpiresAt, err := helper.GenerateToken(dbUser.User_name, accessTokenExpiration)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusInternalServerError)
	}

	refreshToken, _, err := helper.GenerateToken(dbUser.User_name, refreshTokenExpiration)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusInternalServerError)
	}

	signings := &models.Signings{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiration: time.Unix(accessTokenExpiresAt, 0).Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(signings)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user models.User
	var restaurant []models.Restaurant

	database.DB.First(&user, id)
	database.DB.Model(&user).Related(&restaurant)

	user.Restaurant = restaurant

	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("User does not exist")
		return
	}

	userData := map[string]interface{}{
		"id":           user.ID,
		"first_name":   user.First_name,
		"last_name":    user.Last_name,
		"user_name":    user.User_name,
		"email":        user.Email,
		"avatar":       user.Avatar,
		"phone":        user.Phone,
		"account_type": user.Account_type,
		"restaurant":   user.Restaurant,
	}

	json.NewEncoder(w).Encode(userData)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user models.User
	var dbUser models.User
	var restaurant []models.Restaurant

	database.DB.First(&dbUser, id)
	database.DB.Model(&dbUser).Related(&restaurant)

	dbUser.Restaurant = restaurant

	if dbUser.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("User does not exist")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&user)

	file, _, _ := r.FormFile("avatar")
	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "avatar")
		if err != nil {
			dbUser.Avatar = user.Avatar
		}
		dbUser.Avatar = avatarUrl
	}

	dbUser.First_name = user.First_name
	dbUser.Last_name = user.Last_name

	database.DB.Save(&dbUser)

	userData := map[string]interface{}{
		"id":           dbUser.ID,
		"first_name":   dbUser.First_name,
		"last_name":    dbUser.Last_name,
		"user_name":    dbUser.User_name,
		"email":        dbUser.Email,
		"avatar":       dbUser.Avatar,
		"phone":        dbUser.Phone,
		"account_type": dbUser.Account_type,
		"restaurant":   dbUser.Restaurant,
	}

	json.NewEncoder(w).Encode(userData)
}
