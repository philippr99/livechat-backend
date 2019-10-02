package routes

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/http"
	"symflower/livechat_gqlgen/auth"
	"symflower/livechat_gqlgen/management"
	"symflower/livechat_gqlgen/models"
)

func CreateAccountEndPoint() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var user models.Account
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			return
		}

		if len(user.Username) < 4 {
			return
		}

		if management.IsUserAlreadyRegisterd(user) {
			fmt.Printf("[Error] User [%s] is already registered.\n", user.Username)
			return
		}

		fmt.Printf("User [%s] is logging in.\n", user.Username)
		userUuid, _ := uuid.NewV4()
		management.AddUser(user)
		user.ID = userUuid.String()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       user.ID,
			"username": user.Username,
		})
		tokenString, error := token.SignedString(auth.JWTSecret)

		if error != nil {
			log.Println(err)
			return
		}
		response.Header().Set("content-type", "application/json")
		response.Write([]byte(`{"token": "` + tokenString + `"}`))
	}
}
