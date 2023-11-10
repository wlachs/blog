package services

import (
	"fmt"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/repository"
	"github.com/wlchs/blog/internal/types"
	"os"
)

// UserService interface. Defines user-related business logic.
type UserService interface {
}

// userService is the concrete implementation of the UserService interface.
type userService struct {
	cont container.Container
}

// GetUser retrieves a user by userName and creates a user data object.
func (u userService) GetUser(userName string) (types.User, error) {
	userRepository := u.cont.GetUserRepository()
	user, err := userRepository.GetUser(userName)
	return mapUser(user), err
}

// GetUsers retrieves every user and maps them to a slice of user data objects.
func (u userService) GetUsers() ([]types.User, error) {
	userRepository := u.cont.GetUserRepository()
	users, err := userRepository.GetUsers()
	return mapUsers(users), err
}

// RegisterFirstUser creates the main user if it doesn't exist yet.
// The default username and password are read from environment variables.
func (u userService) RegisterFirstUser() error {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	defaultUser := os.Getenv("DEFAULT_USER")
	defaultPassword := os.Getenv("DEFAULT_PASSWORD")

	if defaultUser == "" || defaultPassword == "" {
		return fmt.Errorf("default username or password missing")
	}

	if _, userNotFound := userRepository.GetUser(defaultUser); userNotFound == nil {
		log.Infof("default user with name %s already exists", defaultPassword)
		return nil
	}

	user := types.UserLoginInput{
		UserName: defaultUser,
		Password: defaultPassword,
	}

	log.Infof("initializing first user with name %s", defaultUser)
	_, err := u.RegisterUser(&user)
	return err
}

// RegisterUser creates a new user with the provided username and password.
func (u userService) RegisterUser(user *types.UserLoginInput) (types.User, error) {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	hash, err := HashString(user.Password)
	if err != nil {
		log.Debugf("failed to calculate password hash: %v", err)
		return types.User{}, err
	}

	newUser := types.User{
		UserName:     user.UserName,
		PasswordHash: hash,
	}

	addedUser, err := userRepository.AddUser(&newUser)
	return mapUser(addedUser), err
}

// UpdateUser receives two user input objects, one with the user's current password, and one with the new attributes.
// If the old password matches the currently set one, the new fields are set.
func (u userService) UpdateUser(oldUser *types.UserLoginInput, newUser *types.UserLoginInput) (types.User, error) {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	if ok := CheckUserPassword(*oldUser); !ok {
		log.Debugf("incorrect password for user: %s", oldUser.UserName)
		return types.User{}, errortypes.IncorrectUsernameOrPasswordError{}
	}

	hash, err := HashString(newUser.Password)
	if err != nil {
		log.Debugf("failed to hash new password for user: %s", newUser.UserName)
		return types.User{}, errortypes.PasswordHashingError{}
	}

	user := types.User{
		UserName:     newUser.UserName,
		PasswordHash: hash,
	}

	updatedUser, err := userRepository.UpdateUser(&user)
	if err != nil {
		log.Debugf("failed to update user: %s", user.UserName)
		return types.User{}, err
	}

	log.Debugf("updated user: %s", user.UserName)
	return mapUser(updatedUser), nil
}

// mapUSer maps a User model to a user data object
func mapUser(u *repository.User) types.User {
	return types.User{
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
		Posts:        mapPostHandles(u.Posts),
	}
}

// mapUsers maps a slice of User models to a slice of user data objects
func mapUsers(u []repository.User) []types.User {
	users := make([]types.User, 0, len(u))

	for _, user := range u {
		users = append(users, mapUser(&user))
	}

	return users
}
