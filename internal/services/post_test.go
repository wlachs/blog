package services_test

import (
	"github.com/golang/mock/gomock"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/services"
	"testing"
)

// postTestContext contains objects relevant for testing the PostService.
type postTestContext struct {
	sut services.PostService
}

// createPostServiceContext creates the context for testing the PostService and reduces code duplication.
func createPostServiceContext(t *testing.T) *postTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockPostRepository := mocks.NewMockPostRepository(mockCtrl)
	mockUserRepository := mocks.NewMockUserRepository(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), mockPostRepository, mockUserRepository, nil)
	sut := services.CreatePostService(cont)

	return &postTestContext{sut}
}
