package services

import (
	"os"

	"github.com/wlchs/blog/internal/errors"
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/types"
)

func mapUser(u models.User) types.User {
	return types.User{
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
		Posts:        mapPostHandles(u.Posts),
	}
}

func mapUsers(u []models.User) []types.User {
	users := []types.User{}

	for _, user := range u {
		users = append(users, mapUser(user))
	}

	return users
}

func GetUsers() ([]types.User, error) {
	u, err := models.GetUsers()
	return mapUsers(u), err
}

func GetUser(userName string) (types.User, error) {
	u, err := models.GetUser(userName)
	return mapUser(u), err
}

func RegisterUser(u types.UserRegisterInput) (types.User, error) {
	REGISTRATION_SECRET := os.Getenv("REGISTRATION_SECRET")

	if REGISTRATION_SECRET != u.RegistrationSecret {
		return types.User{}, errors.IncorrectSecretError{}
	}

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
