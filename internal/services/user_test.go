package services_test

import (
	"github.com/golang/mock/gomock"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/services"
	"testing"
)

// userTestContext contains objects relevant for testing the UserService.
type userTestContext struct {
	sut services.UserService
}

// createUserServiceContext creates the context for testing the UserService and reduces code duplication.
func createUserServiceContext(t *testing.T) *userTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockPostRepository := mocks.NewMockPostRepository(mockCtrl)
	mockUserRepository := mocks.NewMockUserRepository(mockCtrl)
	mockJwtUtils := mocks.NewMockTokenUtils(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), mockPostRepository, mockUserRepository, mockJwtUtils)
	sut := services.CreateUserService(cont)

	return &userTestContext{sut}
}
