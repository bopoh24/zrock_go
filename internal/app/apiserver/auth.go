package apiserver

import (
	"time"

	"github.com/bopoh24/zrock_go/internal/app/settings"
	jwt "github.com/dgrijalva/jwt-go"
)

type tokenClaims struct {
	UserID int
	jwt.StandardClaims
}

// CreateJWT creates JWT
func CreateJWT(id int) (string, error) {
	claims := tokenClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(settings.App.TokenSecret))
}

// ParseJWT parses JWT and returns user id and error
func ParseJWT(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.App.TokenSecret), nil
	})
	if err != nil {
		return 0, err
	}
	if claim, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claim.UserID, err
	}
	return 0, err
}
