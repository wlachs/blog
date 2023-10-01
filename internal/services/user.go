package services

import (
	"fmt"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/types"
	"github.com/wlchs/blog/internal/utils"
	"net/http"
	"os"
)

func mapUser(u models.User) types.User {
	return types.User{
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
		Posts:        mapPostHandles(u.Posts),
	}
}

func mapUsers(u []models.User) []types.User {
	var users []types.User

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

// RegisterFirstUser creates the main user if it doesn't exist yet. The default username and password are read from environment variables.
func RegisterFirstUser() error {
	defaultUser := os.Getenv("DEFAULT_USER")
	defaultPassword := os.Getenv("DEFAULT_PASSWORD")

	if defaultUser == "" || defaultPassword == "" {
		return fmt.Errorf("default username or password missing")
	}

	if _, userNotFound := GetUser(defaultUser); userNotFound == nil {
		return nil
	}

	u := types.UserLoginInput{
		UserName: defaultUser,
		Password: defaultPassword,
	}

	utils.LOG.Infof("initializing first user with name %s", defaultUser)
	_, err := RegisterUser(u)
	return err
}

func RegisterUser(u types.UserLoginInput) (types.User, error) {
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

// ChangeUserPassword receives two user input objects, one with the user's current password, and one with the new one.
// If the old password matches the currently set one, the new password is hashed and set instead.
func ChangeUserPassword(oldUser types.UserLoginInput, newUser types.UserLoginInput) (types.User, error) {
	if ok := CheckUserPassword(oldUser); !ok {
		return types.User{}, errortypes.ErrorWithStatus{Status: 401, Err: fmt.Errorf("incorrect username or password")}
	}

	hash, err := HashString(newUser.Password)
	if err != nil {
		return types.User{}, errortypes.ErrorWithStatus{Status: http.StatusUnauthorized, Err: err}
	}

	user := types.User{UserName: oldUser.UserName, PasswordHash: hash}
	updatedUser, err := models.UpdateUser(user)

	return mapUser(updatedUser), nil
}
