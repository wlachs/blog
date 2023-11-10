package repository

import (
	"fmt"
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
}

// userRepository is the concrete implementation of the UserRepository interface
type userRepository struct {
	logger     *zap.SugaredLogger
	repository Repository
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
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		log.Debugf("no user found: %v", user)
		return nil, fmt.Errorf("user with name: %s not found", userName)
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
		log.Debug("failed to retrieve users: %v", result.Error)
		return []User{}, result.Error
	}

	log.Debug("retrieved users: %v", users)
	return users, nil
}

// UpdateUser updates an existing user with the provided data.
func (u userRepository) UpdateUser(user types.User) (*User, error) {
	log := u.logger
	repo := u.repository

	userToUpdate := User{UserName: user.UserName}
	pw := User{PasswordHash: user.PasswordHash}

	if result := repo.Where(&user).Updates(&pw); result.Error != nil {
		log.Debugf("failed to update user %v, error: %v", userToUpdate, result.Error)
		return nil, result.Error
	}

	log.Debugf("updated user: %v", userToUpdate)
	return &userToUpdate, nil
}
