package management

import "symflower/livechat_gqlgen/models"

var registeredUsers = []models.Account{}

func AddUser(user models.Account) {
	registeredUsers = append(registeredUsers, user)
}

func IsUserAlreadyRegisterd(user models.Account) (bool) {
	for _, curUser:= range registeredUsers {
		if curUser.Username == user.Username {
			return true
		}
	}
	return false
}
