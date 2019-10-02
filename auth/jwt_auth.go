package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"symflower/livechat_gqlgen/models"
)

var JWTSecret []byte = []byte("livechatSecret")

func ValidateJWT(tokenString string) (interface{}, error) {
	if tokenString == "" {
		return nil, errors.New("No authentification token provided")
	}

	parsedToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error while parsing provided token")
		}
		return JWTSecret, nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		var decodedToken models.Account
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	} else {
		return nil, errors.New("Invalid authentification token provided")
	}
}