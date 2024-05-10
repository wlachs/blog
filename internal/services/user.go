package services

//go:generate mockgen-v0.4.0 -source=user.go -destination=../mocks/mock_user_service.go -package=mocks

import (
	"github.com/wlachs/blog/internal/auth"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/repository"
	"github.com/wlachs/blog/internal/types"
	"os"
)

// UserService interface. Defines user-related business logic.
type UserService interface {
	AuthenticateUser(user *types.UserLoginInput) (string, error)
	CheckUserPassword(user *types.UserLoginInput) bool
	GetUser(userName string) (types.User, error)
	GetUsers() ([]types.User, error)
	RegisterFirstUser() error
	RegisterUser(user *types.UserLoginInput) (types.User, error)
	UpdateUser(oldUser *types.UserLoginInput, newUser *types.UserLoginInput) (types.User, error)
}

// userService is the concrete implementation of the UserService interface.
type userService struct {
	cont container.Container
}

// CreateUserService instantiates the userService using the application container.
func CreateUserService(cont container.Container) UserService {
	u := &userService{cont}
	initUserService(u)
	return u
}

// initUserService contains logic that should be executed directly upon initialization.
// For now, it takes care of adding the main user to the system if it doesn't exist.
func initUserService(service *userService) {
	log := service.cont.GetLogger()

	// Create user if it doesn't exist yet
	if err := service.RegisterFirstUser(); err != nil {
		log.Errorf("first user registration failed: %s", err)
		return
	}

	log.Infoln("init actions done")
}

// AuthenticateUser authenticates the user.
// If the password hash matches the one stored in the database, a JWT is generated.
func (u userService) AuthenticateUser(user *types.UserLoginInput) (string, error) {
	log := u.cont.GetLogger()
	jwtUtils := u.cont.GetJWTUtils()

	if !u.CheckUserPassword(user) {
		log.Debugf("the provided password hash for user \"%s\" doesn't match the one stored in the DB", user.UserName)
		return "", errortypes.IncorrectUsernameOrPasswordError{}
	}

	log.Debugf("authentication complete for user: %s", user.UserName)
	return jwtUtils.GenerateJWT(user.UserName)
}

// CheckUserPassword fetches the user's password hash from the database and compares it to the input.
func (u userService) CheckUserPassword(user *types.UserLoginInput) bool {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	userModel, err := userRepository.GetUser(user.UserName)
	if err != nil {
		log.Errorf("failed to get user %s from DB: %v", user.UserName, err)
		return false
	}

	return auth.CompareStringWithHash(user.Password, userModel.PasswordHash)
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
		return errortypes.MissingDefaultUsernameOrPasswordError{}
	}

	if _, userNotFound := userRepository.GetUser(defaultUser); userNotFound == nil {
		log.Infof("default user with name %s already exists", defaultUser)
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

	hash, err := auth.HashString(user.Password)
	if err != nil {
		log.Errorf("failed to calculate password hash: %v", err)
		return types.User{}, errortypes.PasswordHashingError{}
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

	if ok := u.CheckUserPassword(oldUser); !ok {
		log.Debugf("incorrect password for user: %s", oldUser.UserName)
		return types.User{}, errortypes.IncorrectUsernameOrPasswordError{}
	}

	hash, err := auth.HashString(newUser.Password)
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
	if u == nil {
		return types.User{}
	}
	return types.User{
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
		Posts:        mapPostHandles(u.Posts),
	}
}

// mapUsers maps a slice of User models to a slice of user data objects
func mapUsers(u []repository.User) []types.User {
	if u == nil {
		return []types.User{}
	}
	users := make([]types.User, 0, len(u))

	for _, user := range u {
		users = append(users, mapUser(&user))
	}

	return users
}
