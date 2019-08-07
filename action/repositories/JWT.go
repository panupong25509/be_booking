package repositories

import (
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/panupong25509/be_booking_sign/models"
)

type Token struct {
	UserID uuid.UUID
	Role   string
	jwt.StandardClaims
}

func EncodeJWT(user models.User, secret string) string {
	tokenJWT := Token{UserID: user.ID, Role: user.Role}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenJWT)
	jwt, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalln(err)
	}
	return jwt
}

func DecodeJWT(jwtReq string, secret string) (jwt.MapClaims, interface{}) {
	jwtStrings := strings.Split(jwtReq, "Bearer ")
	token, err := jwt.Parse(jwtStrings[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, models.Error{500, "Token error."}
	}
	tokens := token.Claims.(jwt.MapClaims)
	return tokens, nil
}
