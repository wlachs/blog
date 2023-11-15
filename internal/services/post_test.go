package services_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/repository"
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/types"
	"testing"
	"time"
)

// postTestContext contains objects relevant for testing the PostService.
type postTestContext struct {
	mostPostRepository *mocks.MockPostRepository
	mostUserRepository *mocks.MockUserRepository
	sut                services.PostService
}

// createPostServiceContext creates the context for testing the PostService and reduces code duplication.
func createPostServiceContext(t *testing.T) *postTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockPostRepository := mocks.NewMockPostRepository(mockCtrl)
	mockUserRepository := mocks.NewMockUserRepository(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), mockPostRepository, mockUserRepository, nil)
	sut := services.CreatePostService(cont)

	return &postTestContext{mockPostRepository, mockUserRepository, sut}
}

// TestPostService_AddPost tests adding a new post to the blog.
func TestPostService_AddPost(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	userModel := repository.User{
		ID:       0,
		UserName: "testAuthor",
		Posts:    []repository.Post{},
	}

	postModel := repository.Post{
		ID:        0,
		URLHandle: "testUrlHandle",
		AuthorID:  userModel.ID,
		Author:    userModel,
		Title:     "testTitle",
		Summary:   "testSummary",
		Body:      "testBody",
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}

	newPost := types.Post{
		URLHandle:    postModel.URLHandle,
		Title:        postModel.Title,
		Author:       userModel.UserName,
		Summary:      postModel.Summary,
		Body:         postModel.Body,
		CreationTime: postModel.CreatedAt,
	}

	c.mostUserRepository.EXPECT().GetUser(userModel.UserName).Return(&userModel, nil)
	c.mostPostRepository.EXPECT().AddPost(&newPost, &userModel).Return(&postModel, nil)

	p, err := c.sut.AddPost(&newPost)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, newPost, p, "added post doesn't match the input")
}

// TestPostService_AddPost_Invalid_User tests adding a new post to the blog with invalid username.
func TestPostService_AddPost_Invalid_User(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	newPost := types.Post{
		URLHandle: "testUrlHandle",
		Title:     "testTitle",
		Author:    "testAuthor",
		Summary:   "testSummary",
		Body:      "testBody",
	}

	c.mostUserRepository.EXPECT().GetUser("testAuthor").Return(nil, fmt.Errorf("error"))

	p, err := c.sut.AddPost(&newPost)

	assert.NotNil(t, err, "expected error")
	assert.NotEqual(t, newPost, p, "added user with incorrect data")
}

// TestPostService_Duplicate_Post tests adding a duplicate post to the blog.
func TestPostService_Duplicate_Post(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	userModel := repository.User{
		ID:       0,
		UserName: "testAuthor",
		Posts:    []repository.Post{},
	}

	postModel := repository.Post{
		ID:        0,
		URLHandle: "duplicateUrlHandle",
		AuthorID:  userModel.ID,
		Author:    userModel,
		Title:     "testTitle",
		Summary:   "testSummary",
		Body:      "testBody",
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}

	newPost := types.Post{
		URLHandle:    postModel.URLHandle,
		Title:        postModel.Title,
		Author:       userModel.UserName,
		Summary:      postModel.Summary,
		Body:         postModel.Body,
		CreationTime: postModel.CreatedAt,
	}

	expectedError := errortypes.DuplicateElementError{Key: postModel.URLHandle}

	c.mostUserRepository.EXPECT().GetUser(userModel.UserName).Return(&userModel, nil)
	c.mostPostRepository.EXPECT().AddPost(&newPost, &userModel).Return(nil, expectedError)

	p, err := c.sut.AddPost(&newPost)

	assert.Equal(t, expectedError, err, "error doesn't match expected one")
	assert.NotEqual(t, newPost, p, "added user with incorrect data")
}

// TestPostService_GetPost tests getting a post from the blog.
func TestPostService_GetPost(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	userModel := repository.User{
		ID:       0,
		UserName: "testAuthor",
		Posts:    []repository.Post{},
	}

	postModel := repository.Post{
		ID:        0,
		URLHandle: "testUrlHandle",
		AuthorID:  userModel.ID,
		Author:    userModel,
		Title:     "testTitle",
		Summary:   "testSummary",
		Body:      "testBody",
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}

	post := types.Post{
		URLHandle:    postModel.URLHandle,
		Title:        postModel.Title,
		Author:       userModel.UserName,
		Summary:      postModel.Summary,
		Body:         postModel.Body,
		CreationTime: postModel.CreatedAt,
	}

	c.mostPostRepository.EXPECT().GetPost(postModel.URLHandle).Return(&postModel, nil)

	p, err := c.sut.GetPost(postModel.URLHandle)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, post, p, "post doesn't match the expected output")
}

// TestPostService_GetPost_Unexpected_Error tests handling an unexpected error while getting a post.
func TestPostService_GetPost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	c.mostPostRepository.EXPECT().GetPost("testUrlHandle").Return(nil, fmt.Errorf("error"))
	_, err := c.sut.GetPost("testUrlHandle")

	assert.NotNil(t, err, "expected error")
}

// TestPostService_GetPosts tests getting posts from the blog.
func TestPostService_GetPosts(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	userModel := repository.User{
		ID:       0,
		UserName: "testAuthor",
		Posts:    []repository.Post{},
	}

	postModels := []repository.Post{
		{
			ID:        0,
			URLHandle: "testUrlHandle",
			AuthorID:  userModel.ID,
			Author:    userModel,
			Title:     "testTitle",
			Summary:   "testSummary",
			Body:      "testBody",
			CreatedAt: time.Time{}.Local(),
			UpdatedAt: time.Time{}.Local(),
		},
	}

	posts := []types.Post{
		{
			URLHandle:    postModels[0].URLHandle,
			Title:        postModels[0].Title,
			Author:       userModel.UserName,
			Summary:      postModels[0].Summary,
			CreationTime: postModels[0].CreatedAt,
		},
	}

	c.mostPostRepository.EXPECT().GetPosts().Return(postModels, nil)

	p, err := c.sut.GetPosts()

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, posts, p, "post doesn't match the expected output")
}

// TestPostService_GetPosts_Unexpected_Error tests handling an unexpected error while getting posts
func TestPostService_GetPosts_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	c.mostPostRepository.EXPECT().GetPosts().Return(nil, fmt.Errorf("error"))
	_, err := c.sut.GetPosts()

	assert.NotNil(t, err, "expected error")
}
