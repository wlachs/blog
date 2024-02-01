package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
