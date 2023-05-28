package common

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenToken(id primitive.ObjectID) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))

	jwt_token.Claims = jwt.MapClaims{
		"id": id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	jwtSecret := GetENV("JWT_SECRET")

	token, _ := jwt_token.SignedString([]byte(jwtSecret))
	return token
}

func NormalizeString(str string) string {
    str = strings.ToLower(str)
    str = strings.ReplaceAll(str, " ", "")
    return str
}