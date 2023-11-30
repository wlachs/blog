package repository

import (
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/types"
	"go.uber.org/zap"
	"time"
)

// User DB schema
type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	UserName     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Posts        []Post `gorm:"foreignKey:AuthorID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository interface defining user-related database operations.
type UserRepository interface {
	AddUser(user *types.User) (*User, error)
	GetUser(userName string) (*User, error)
	GetUsers() ([]User, error)
	UpdateUser(user *types.User) (*User, error)
}

// userRepository is the concrete implementation of the UserRepository interface
type userRepository struct {
	logger     *zap.SugaredLogger
	repository Repository
}

// CreateUserRepository instantiates the userRepository using the logger and the global repository.
func CreateUserRepository(logger *zap.SugaredLogger, repository Repository) UserRepository {
	initUserModel(logger, repository)

	return &userRepository{
		logger:     logger,
		repository: repository,
	}
}

// initUserModel initializes the User schema in the database
func initUserModel(logger *zap.SugaredLogger, repository Repository) {
	if err := repository.AutoMigrate(&User{}); err != nil {
		logger.Errorf("failed to initialize user model: %v", err)
	}
}

// AddUser adds a new user with the provided fields to the database.
func (u userRepository) AddUser(user *types.User) (*User, error) {
	log := u.logger
	repo := u.repository

	newUser := User{
		UserName:     user.UserName,
		PasswordHash: user.PasswordHash,
	}

	if result := repo.Create(&newUser); result.Error != nil {
		log.Debugf("failed to create new user: %v, error: %v", newUser, result.Error)
		return nil, result.Error
	}

	log.Debugf("created new user: %v", newUser)
	return &newUser, nil
}

// GetUser retrieves a user with the given userName from the database.
func (u userRepository) GetUser(userName string) (*User, error) {
	log := u.logger
	repo := u.repository

	user := User{
		UserName: userName,
	}

	result := repo.Preload("Posts").Where(&user).Take(&user)

	if result.Error != nil {
		log.Debugf("failed to retrieve user: %v, error: %v", user, result.Error)
		if result.Error.Error() == "record not found" {
			return nil, errortypes.UserNotFoundError{User: types.User{UserName: userName}}
		}
		return nil, result.Error
	}

	log.Debugf("retrieved user: %v", user)
	return &user, nil
}

// GetUsers retrieves every user from the database.
func (u userRepository) GetUsers() ([]User, error) {
	log := u.logger
	repo := u.repository

	var users []User
	if result := repo.Preload("Posts").Find(&users); result.Error != nil {
		log.Debugf("failed to retrieve users: %v", result.Error)
		return []User{}, result.Error
	}

	log.Debugf("retrieved users: %v", users)
	return users, nil
}

// UpdateUser updates an existing user with the provided data.
func (u userRepository) UpdateUser(user *types.User) (*User, error) {
	log := u.logger
	repo := u.repository

	userToUpdate := User{UserName: user.UserName}
	pw := User{PasswordHash: user.PasswordHash}

	if result := repo.Where(&userToUpdate).Updates(&pw); result.Error != nil {
		log.Debugf("failed to update user %v, error: %v", userToUpdate, result.Error)
		return nil, result.Error
	}

	log.Debugf("updated user: %v", userToUpdate)
	return &userToUpdate, nil
}
