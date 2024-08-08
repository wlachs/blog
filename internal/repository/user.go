package repository

//go:generate mockgen-v0.4.0 -source=user.go -destination=../mocks/mock_user_repository.go -package=mocks

import (
	"github.com/wlachs/blog/internal/errortypes"
	"go.uber.org/zap"
	"strings"
	"time"
)

// User DB schema
type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	UserName     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Posts        []Post `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository interface defining user-related database operations.
type UserRepository interface {
	AddUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(userName string) error
	GetUser(userName string) (User, error)
	GetUsers(pageIndex int, pageSize int) ([]User, error)
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
func (u userRepository) AddUser(user User) (User, error) {
	log := u.logger
	repo := u.repository

	if result := repo.Create(&user); result.Error != nil {
		if strings.Contains(result.Error.Error(), "1062") {
			log.Debugf("failed to create new user, duplicate key: %s, error: %v", user.UserName, result.Error)
			return User{}, errortypes.DuplicateElementError{Key: user.UserName}
		} else {
			log.Debugf("failed to create new user: %v, error: %v", user, result.Error)
			return User{}, result.Error
		}
	}

	populateUserAsAuthorOfPosts(&user)

	log.Debugf("created new user: %v", user)
	return user, nil
}

// UpdateUser updates an existing user with the provided data.
func (u userRepository) UpdateUser(user User) (User, error) {
	log := u.logger
	repo := u.repository

	userToUpdate := User{UserName: user.UserName}
	pw := User{PasswordHash: user.PasswordHash}

	if result := repo.Where(&userToUpdate).Updates(&pw); result.Error != nil {
		log.Debugf("failed to update user %v, error: %v", userToUpdate, result.Error)
		return User{}, result.Error
	}

	populateUserAsAuthorOfPosts(&user)

	log.Debugf("updated user: %v", userToUpdate)
	return userToUpdate, nil
}

// DeleteUser deletes a user from the database.
func (u userRepository) DeleteUser(userName string) error {
	log := u.logger
	repo := u.repository

	user := User{UserName: userName}

	if result := repo.Where(user).Delete(user); result.Error == nil {
		if result.RowsAffected > 0 {
			log.Debugf("deleted user: %s", userName)
			return nil
		} else {
			log.Debugf("failed to delete user: %v, user not found", user)
			return errortypes.UserNotFoundError{UserName: userName}
		}
	} else {
		log.Debugf("failed to delete user: %v, error: %s", user, result.Error)
		return result.Error
	}
}

// GetUser retrieves a user with the given userName from the database.
func (u userRepository) GetUser(userName string) (User, error) {
	log := u.logger
	repo := u.repository

	user := User{
		UserName: userName,
	}

	result := repo.Preload("Posts").Where(&user).Take(&user)

	if result.Error != nil {
		log.Debugf("failed to retrieve user: %v, error: %v", user, result.Error)
		if result.Error.Error() == "record not found" {
			return User{}, errortypes.UserNotFoundError{UserName: userName}
		}
		return User{}, result.Error
	}

	populateUserAsAuthorOfPosts(&user)

	log.Debugf("retrieved user: %v", user)
	return user, nil
}

// GetUsers retrieves a specific page of users from the database.
func (u userRepository) GetUsers(pageIndex int, pageSize int) ([]User, error) {
	log := u.logger
	repo := u.repository

	var users []User
	result := repo.Preload("Posts").
		Limit(pageSize).
		Offset((pageIndex - 1) * pageSize).
		Find(&users)

	if result.Error != nil {
		log.Debugf("failed to retrieve users: %v", result.Error)
		return []User{}, result.Error
	}

	populateUsersAsAuthorsOfPosts(users)

	log.Debugf("retrieved users: %v", users)
	return users, nil
}

// populateUserAsAuthorOfPosts manually sets user model for contained posts
func populateUserAsAuthorOfPosts(user *User) {
	for i := range user.Posts {
		user.Posts[i].Author = *user
		// Set posts to nil to avoid circular reference
		user.Posts[i].Author.Posts = nil
	}
}

// populateUsersAsAuthorsOfPosts manually sets user models for contained posts
func populateUsersAsAuthorsOfPosts(users []User) {
	for userIndex := range users {
		populateUserAsAuthorOfPosts(&users[userIndex])
	}
}
