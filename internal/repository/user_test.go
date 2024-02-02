package repository_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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
