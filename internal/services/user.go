package services

//go:generate mockgen-v0.4.0 -source=user.go -destination=../mocks/mock_user_service.go -package=mocks

import (
	"github.com/wlachs/blog/internal/auth"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/repository"
	"os"
)

// UserService interface. Defines user-related business logic.
type UserService interface {
	AuthenticateUser(userID string, password string) (string, error)
	CheckUserPassword(userID string, password string) bool
	GetUser(userName string) (repository.User, error)
	GetUsers() ([]repository.User, error)
	RegisterFirstUser() error
	RegisterUser(userID string, password string) (repository.User, error)
	UpdateUser(userID string, oldPassword string, newPassword string) (repository.User, error)
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
func (u userService) AuthenticateUser(userID string, password string) (string, error) {
	log := u.cont.GetLogger()
	jwtUtils := u.cont.GetJWTUtils()

	if !u.CheckUserPassword(userID, password) {
		log.Debugf("the provided password hash for user \"%s\" doesn't match the one stored in the DB", userID)
		return "", errortypes.IncorrectUsernameOrPasswordError{}
	}

	log.Debugf("authentication complete for user: %s", userID)
	return jwtUtils.GenerateJWT(userID)
}

// CheckUserPassword fetches the user's password hash from the database and compares it to the input.
func (u userService) CheckUserPassword(userID string, password string) bool {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	userModel, err := userRepository.GetUser(userID)
	if err != nil {
		log.Errorf("failed to get user %s from DB: %v", userID, err)
		return false
	}

	return auth.CompareStringWithHash(password, userModel.PasswordHash)
}

// GetUser retrieves a user by userName and creates a user data object.
func (u userService) GetUser(userID string) (repository.User, error) {
	userRepository := u.cont.GetUserRepository()
	return userRepository.GetUser(userID)
}

// GetUsers retrieves every user and maps them to a slice of user data objects.
func (u userService) GetUsers() ([]repository.User, error) {
	userRepository := u.cont.GetUserRepository()
	return userRepository.GetUsers()
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

	log.Infof("initializing first user with name %s", defaultUser)
	_, err := u.RegisterUser(defaultUser, defaultPassword)
	return err
}

// RegisterUser creates a new user with the provided username and password.
func (u userService) RegisterUser(userID string, password string) (repository.User, error) {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	hash, err := auth.HashString(password)
	if err != nil {
		log.Errorf("failed to calculate password hash: %v", err)
		return repository.User{}, errortypes.PasswordHashingError{}
	}

	newUser := repository.User{
		UserName:     userID,
		PasswordHash: hash,
	}

	return userRepository.AddUser(newUser)
}

// UpdateUser receives two user input objects, one with the user's current password, and one with the new attributes.
// If the old password matches the currently set one, the new fields are set.
func (u userService) UpdateUser(userID string, oldPassword string, newPassword string) (repository.User, error) {
	log := u.cont.GetLogger()
	userRepository := u.cont.GetUserRepository()

	if ok := u.CheckUserPassword(userID, oldPassword); !ok {
		log.Debugf("incorrect password for user: %s", userID)
		return repository.User{}, errortypes.IncorrectUsernameOrPasswordError{}
	}

	hash, err := auth.HashString(newPassword)
	if err != nil {
		log.Debugf("failed to hash new password for user: %s", userID)
		return repository.User{}, errortypes.PasswordHashingError{}
	}

	user := repository.User{
		UserName:     userID,
		PasswordHash: hash,
	}

	updatedUser, err := userRepository.UpdateUser(user)
	if err != nil {
		log.Debugf("failed to update user: %s", user.UserName)
		return repository.User{}, err
	}

	log.Debugf("updated user: %s", user.UserName)
	return updatedUser, nil
}
