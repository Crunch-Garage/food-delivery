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
	json.NewEncoder(w).Encode(user_.Profile)
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

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id, _ := strconv.Atoi(params["id"])

// 	var user models.User
// 	var dbUser models.User
// 	var restaurant []models.Restaurant

// 	database.DB.First(&dbUser, id)
// 	database.DB.Model(&dbUser).Related(&restaurant)

// 	dbUser.Restaurant = restaurant

// 	if dbUser.ID == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("User does not exist")
// 		return
// 	}

// 	_ = json.NewDecoder(r.Body).Decode(&user)

// 	file, _, _ := r.FormFile("avatar")
// 	if file != nil {
// 		avatarUrl, err := helper.SingleImageUpload(w, r, "avatar")
// 		if err != nil {
// 			dbUser.Avatar = user.Avatar
// 		}
// 		dbUser.Avatar = avatarUrl
// 	}

// 	dbUser.First_name = user.First_name
// 	dbUser.Last_name = user.Last_name

// 	database.DB.Save(&dbUser)

// 	userData := map[string]interface{}{
// 		"id":           dbUser.ID,
// 		"first_name":   dbUser.First_name,
// 		"last_name":    dbUser.Last_name,
// 		"user_name":    dbUser.User_name,
// 		"email":        dbUser.Email,
// 		"avatar":       dbUser.Avatar,
// 		"phone":        dbUser.Phone,
// 		"account_type": dbUser.Account_type,
// 		"restaurant":   dbUser.Restaurant,
// 	}

// 	json.NewEncoder(w).Encode(userData)
// }
