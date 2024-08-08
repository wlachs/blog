package repository_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

// userTestContext contains objects relevant for testing the UserRepository.
type userTestContext struct {
	mockDb sqlmock.Sqlmock
	sut    repository.UserRepository
}

// createUserRepositoryContext creates the context for testing the UserRepository and reduces code duplication.
func createUserRepositoryContext(t *testing.T) *userTestContext {
	t.Helper()

	db, mock, _ := sqlmock.New()
	gormDb, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))

	sut := repository.CreateUserRepository(logger.CreateLogger(), repository.CreateRepository(gormDb))
	return &userTestContext{mock, sut}
}

// TestUserRepository_AddUser tests adding a new user to the system.
func TestUserRepository_AddUser(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := repository.User{
		UserName: "testUser",
	}

	userQuery := regexp.QuoteMeta("INSERT INTO `users` (`user_name`,`password_hash`,`created_at`,`updated_at`) VALUES (?,?,?,?)")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	c.mockDb.ExpectCommit()

	user, err := c.sut.AddUser(author)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, author.UserName, user.UserName, "received post should match the expected one")
}

// TestUserRepository_AddUser_Duplicate_User tests adding an already existing user to the system.
func TestUserRepository_AddUser_Duplicate_User(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := repository.User{
		UserName: "testUser",
	}

	dbErr := fmt.Errorf("error 1062 (23000): duplicate entry")
	expectedError := errortypes.DuplicateElementError{Key: author.UserName}

	userQuery := regexp.QuoteMeta("INSERT INTO `users` (`user_name`,`password_hash`,`created_at`,`updated_at`) VALUES (?,?,?,?)")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnError(dbErr)
	c.mockDb.ExpectRollback()

	user, err := c.sut.AddUser(author)

	assert.Equal(t, repository.User{}, user, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_AddUser_Unexpected_Error tests adding a new user to the system while encountering an unexpected error
func TestUserRepository_AddUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := repository.User{
		UserName: "testUser",
	}

	expectedError := fmt.Errorf("unexpected error")

	userQuery := regexp.QuoteMeta("INSERT INTO `users` (`user_name`,`password_hash`,`created_at`,`updated_at`) VALUES (?,?,?,?)")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	user, err := c.sut.AddUser(author)

	assert.Equal(t, repository.User{}, user, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_GetUser tests retrieving a single user from the database
func TestUserRepository_GetUser(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	expectedUser := repository.User{
		ID:       0,
		UserName: "testUser",
	}

	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`user_name` = ? LIMIT ?")

	c.mockDb.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name"}).
			AddRow(expectedUser.ID, expectedUser.UserName))

	post, err := c.sut.GetUser(expectedUser.UserName)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, expectedUser, post, "received user should match the expected one")
}

// TestUserRepository_GetUser_Record_Not_Found tests retrieving a non-existent user from the database
func TestUserRepository_GetUser_Record_Not_Found(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	userName := "testUser"

	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`user_name` = ? LIMIT ?")
	dbErr := fmt.Errorf("record not found")
	expectedError := errortypes.UserNotFoundError{UserName: userName}

	c.mockDb.ExpectQuery(query).WillReturnError(dbErr)

	post, err := c.sut.GetUser(userName)

	assert.Equal(t, repository.User{}, post, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_GetUser_Unexpected_Error tests retrieving a single user from the database with an error
func TestUserRepository_GetUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`user_name` = ? LIMIT ?")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectQuery(query).WillReturnError(expectedError)

	post, err := c.sut.GetUser("test")

	assert.Equal(t, repository.User{}, post, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_GetUsers tests retrieving every user from the database
func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	userQuery := regexp.QuoteMeta("SELECT * FROM `users` ORDER BY user_name ASC LIMIT ? OFFSET ?")
	postQuery := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`author_id` IN (?,?)")

	c.mockDb.ExpectQuery(userQuery).
		WithArgs(3, 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name"}).
			AddRow(1, "testUser").
			AddRow(2, "otherTestUser"))

	c.mockDb.ExpectQuery(postQuery).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle", "author_id"}).
			AddRow(1, "test_1", 1).
			AddRow(2, "test_2", 2))

	posts, _, err := c.sut.GetUsers(2, 3)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, 2, len(posts), "didn't receive the expected number of users")
}

// TestUserRepository_GetUsers_Unexpected_Error tests retrieving every user from the database with an error
func TestUserRepository_GetUsers_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `users`")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectQuery(query).WillReturnError(expectedError)

	posts, _, err := c.sut.GetUsers(1, 1)

	assert.Equal(t, expectedError, err, "error should match expected value")
	assert.Equal(t, 0, len(posts), "shouldn't receive any users")
}

// TestUserRepository_UpdateUser tests updating an existing user in the system.
func TestUserRepository_UpdateUser(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := repository.User{
		UserName:     "testUser",
		PasswordHash: "xxx",
	}

	userQuery := regexp.QuoteMeta("UPDATE `users` SET `password_hash`=?,`updated_at`=? WHERE `users`.`user_name` = ?")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	c.mockDb.ExpectCommit()

	user, err := c.sut.UpdateUser(author)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, author.UserName, user.UserName, "received post should match the expected one")
}

// TestUserRepository_UpdateUser_Unexpected_Error tests updating an existing user in the system while encountering an error.
func TestUserRepository_UpdateUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := repository.User{
		UserName:     "testUser",
		PasswordHash: "xxx",
	}

	userQuery := regexp.QuoteMeta("UPDATE `users` SET `password_hash`=?,`updated_at`=? WHERE `users`.`user_name` = ?")

	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	user, err := c.sut.UpdateUser(author)

	assert.Equal(t, repository.User{}, user, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_DeleteUser tests deleting a user from the system without errors.
func TestUserRepository_DeleteUser(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	userQuery := regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`user_name` = ?")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	c.mockDb.ExpectCommit()

	err := c.sut.DeleteUser("testUser")

	assert.Nil(t, err, "should complete without error")
}

// TestUserRepository_DeleteUser_Record_Not_Found tests deleting a non-existing user from the system.
func TestUserRepository_DeleteUser_Record_Not_Found(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	userName := "testUser"
	userQuery := regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`user_name` = ?")
	expectedError := errortypes.UserNotFoundError{UserName: userName}

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	c.mockDb.ExpectCommit()

	err := c.sut.DeleteUser(userName)

	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_DeleteUser_Unexpected_Error tests deleting a user from the system while encountering an error.
func TestUserRepository_DeleteUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	userName := "testUser"
	userQuery := regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`user_name` = ?")

	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	err := c.sut.DeleteUser(userName)

	assert.Equal(t, expectedError, err, "received error should match the expected one")
}
