package controller

import (
	"crunchgarage/restaurant-food-delivery/config"
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

var profile_image = ""

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

/*sign up*/
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	/*check if email exists*/
	var dbUser models.User
	database.DB.Where("email = ?", user.Email).First(&dbUser)
	if dbUser.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Email address already exists")
		return
	}

	/*check if phone phone number exists*/
	database.DB.Where("phone = ?", user.Phone).First(&dbUser)
	if dbUser.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Phone number already exists")
		return
	}

	pro_type := ""

	if user.User_type == "CLIENT" {
		pro_type = ""
	} else {
		pro_type = "CHEF"
	}

	username_email_strip := strings.Split(user.Email, "@")[0]

	user_ := models.User{
		First_name: user.First_name,
		Last_name:  user.Last_name,
		User_type:  user.User_type,
		Email:      user.Email,
		Phone:      user.Phone,
		User_name:  username_email_strip,
		Password:   HashPassword(user.Password),
		Profile: []models.Profile{{
			First_name: user.First_name,
			Last_name:  user.Last_name,
			User_type:  user.User_type,
			User_name:  username_email_strip,
			Pro_type:   pro_type,
		}},
	}

	createdUser := database.DB.Create(&user_)
	err = createdUser.Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	_, err := helper.RegisterEmailAccount(user)

	if err != nil {
		fmt.Println(err)
	}

	/*if successful return user profile*/
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user_.Profile[0])
}

/*login*/
func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User
	var dbUser models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	/*find a user with username and see if that user even exists*/
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
	//user controler
	accessToken, accessTokenExpiresAt, err := helper.GenerateToken(dbUser, accessTokenExpiration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	refreshToken, _, err := helper.GenerateToken(dbUser, refreshTokenExpiration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	signings := &models.Signings{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiration: time.Unix(accessTokenExpiresAt, 0).Format(time.RFC3339),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(signings)
}

/*get user info from profile table*/
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id, _ := strconv.Atoi(params["id"])

	var profile models.Profile
	var restaurant []models.Restaurant

	database.DB.Where("user_id = ?", user_id).First(&profile)
	if profile.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("User does not exist")
		return
	}
	database.DB.Model(&profile).Related(&restaurant)

	profile.Restaurant = restaurant

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user models.User
	var profile models.Profile
	var restaurant []models.Restaurant
	var dbUser models.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	database.DB.First(&dbUser, id)
	database.DB.Where("user_id = ?", id).First(&profile)
	database.DB.Model(&profile).Related(&restaurant)

	if dbUser.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("User does not exist")
		return
	}

	//check for media

	file, _, _ := r.FormFile("profile_image")

	if file != nil {
		avatarUrl, err := helper.SingleImageUpload(w, r, "profile_image", config.EnvCloudMenuFolder())
		if err != nil {
			profile_image = dbUser.Profile_image
		}
		profile_image = avatarUrl
	}

	dbUser.First_name = user.First_name
	dbUser.Last_name = user.Last_name
	dbUser.User_name = user.User_name
	dbUser.Profile_image = profile_image

	/*update user*/
	updatedUser := database.DB.Save(&dbUser)
	err := updatedUser.Error

	/*update profile*/
	database.DB.Model(&profile).Updates(models.Profile{
		First_name:    dbUser.First_name,
		Last_name:     dbUser.Last_name,
		User_name:     dbUser.User_name,
		Profile_image: profile_image,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	profile.Restaurant = restaurant
	dbUser.Profile = []models.Profile{profile}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbUser.Profile)

}
