package repository_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
	"github.com/wlchs/blog/internal/types"
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

	author := &types.User{
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

// TestUserRepository_AddUser_Unexpected_Error tests adding a new user to the system while encountering an unexpected error
func TestUserRepository_AddUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := &types.User{
		UserName: "testUser",
	}

	expectedError := fmt.Errorf("unexpected error")

	userQuery := regexp.QuoteMeta("INSERT INTO `users` (`user_name`,`password_hash`,`created_at`,`updated_at`) VALUES (?,?,?,?)")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	user, err := c.sut.AddUser(author)

	assert.Nil(t, user, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_GetUser tests retrieving a single user from the database
func TestUserRepository_GetUser(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	expectedUser := &repository.User{
		ID:       0,
		UserName: "testUser",
	}

	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`user_name` = ? LIMIT 1")

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

	expectedUser := types.User{
		UserName: "testUser",
	}

	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`user_name` = ? LIMIT 1")
	dbErr := fmt.Errorf("record not found")
	expectedError := errortypes.UserNotFoundError{User: expectedUser}

	c.mockDb.ExpectQuery(query).WillReturnError(dbErr)

	post, err := c.sut.GetUser(expectedUser.UserName)

	assert.Nil(t, post, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_GetUser_Unexpected_Error tests retrieving a single user from the database with an error
func TestUserRepository_GetUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`user_name` = ? LIMIT 1")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectQuery(query).WillReturnError(expectedError)

	post, err := c.sut.GetUser("test")

	assert.Nil(t, post, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestUserRepository_GetUsers tests retrieving every user from the database
func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	userQuery := regexp.QuoteMeta("SELECT * FROM `users`")
	postQuery := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`author_id` IN (?,?)")

	c.mockDb.ExpectQuery(userQuery).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name"}).
			AddRow(1, "testUser").
			AddRow(2, "otherTestUser"))

	c.mockDb.ExpectQuery(postQuery).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle", "author_id"}).
			AddRow(1, "test_1", 1).
			AddRow(2, "test_2", 2))

	posts, err := c.sut.GetUsers()

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

	posts, err := c.sut.GetUsers()

	assert.Equal(t, expectedError, err, "error should match expected value")
	assert.Equal(t, 0, len(posts), "shouldn't receive any users")
}

// TestUserRepository_UpdateUser tests updating an existing user in the system.
func TestUserRepository_UpdateUser(t *testing.T) {
	t.Parallel()
	c := createUserRepositoryContext(t)

	author := &types.User{
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

	author := &types.User{
		UserName:     "testUser",
		PasswordHash: "xxx",
	}

	userQuery := regexp.QuoteMeta("UPDATE `users` SET `password_hash`=?,`updated_at`=? WHERE `users`.`user_name` = ?")

	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(userQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	user, err := c.sut.UpdateUser(author)

	assert.Nil(t, user, "should not return a user")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}
