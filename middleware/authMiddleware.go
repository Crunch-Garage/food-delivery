package middleware

import (
	"crunchgarage/restaurant-food-delivery/models"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte(os.Getenv("JWT_KEY"))

/*Validate JWT requests*/
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return jwtkey, nil
			})
		//	json.NewEncoder(w).Encode(token.Claims)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("Invalid Token")
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Bad Request")
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid Token")
			return
		}

		endpoint(w, r)

	})
}
