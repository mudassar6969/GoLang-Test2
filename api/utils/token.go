package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = []byte("MySecretKey")

func EncodeAuthToken(uid uint) (string, error) {
	claims := jwt.MapClaims{}
	claims[KEY_USER_ID] = uid
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString(SECRET_KEY)
}
