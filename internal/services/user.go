package services

import (
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/types"
)

func mapUser(u models.User) types.User {
	return types.User{
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
	}
}

func GetUser(userName string) (types.User, error) {
	u, err := models.GetUser(userName)
	return mapUser(u), err
}

func RegisterUser(u types.UserInput) (types.User, error) {
	hash, err := HashString(u.Password)

	if err != nil {
		return types.User{}, err
	}

	newUser := types.User{
		UserName:     u.UserName,
		PasswordHash: hash,
	}

	addedUser, err := models.AddUser(newUser)
	return mapUser(addedUser), err
}
