package helper

import (
	"crunchgarage/restaurant-food-delivery/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte(os.Getenv("JWT_KEY"))

/*Genearte JWT token*/
func GenerateToken(principal string, duration time.Duration) (string, int64, error) {

	claims := &models.Claims{
		principal,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)

	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil
}
