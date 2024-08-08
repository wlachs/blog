package services_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/mocks"
	"github.com/wlachs/blog/internal/repository"
	"github.com/wlachs/blog/internal/services"
	"go.uber.org/mock/gomock"
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

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
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
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}
	newPost := repository.Post{
		URLHandle: postModel.URLHandle,
		AuthorID:  userModel.ID,
		Title:     postModel.Title,
		Summary:   postModel.Summary,
		Body:      postModel.Body,
	}

	c.mostUserRepository.EXPECT().GetUser(userModel.UserName).Return(userModel, nil)
	c.mostPostRepository.EXPECT().AddPost(newPost).Return(postModel, nil)

	p, err := c.sut.AddPost(newPost, userModel.UserName)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, postModel, p, "added post doesn't match the input")
}

// TestPostService_AddPost_Invalid_User tests adding a new post to the blog with invalid username.
func TestPostService_AddPost_Invalid_User(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
	newPost := repository.Post{
		URLHandle: "testUrlHandle",
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}

	c.mostUserRepository.EXPECT().GetUser("testAuthor").Return(repository.User{}, fmt.Errorf("error"))

	p, err := c.sut.AddPost(newPost, "testAuthor")

	assert.NotNil(t, err, "expected error")
	assert.NotEqual(t, newPost, p, "added user with incorrect data")
}

// TestPostService_AddPost_Duplicate_Post tests adding a duplicate post to the blog.
func TestPostService_AddPost_Duplicate_Post(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
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
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}
	newPost := repository.Post{
		URLHandle: postModel.URLHandle,
		AuthorID:  userModel.ID,
		Title:     postModel.Title,
		Summary:   postModel.Summary,
		Body:      postModel.Body,
	}

	expectedError := errortypes.DuplicateElementError{Key: postModel.URLHandle}

	c.mostUserRepository.EXPECT().GetUser(userModel.UserName).Return(userModel, nil)
	c.mostPostRepository.EXPECT().AddPost(newPost).Return(repository.Post{}, expectedError)

	p, err := c.sut.AddPost(newPost, userModel.UserName)

	assert.Equal(t, expectedError, err, "error doesn't match expected one")
	assert.NotEqual(t, newPost, p, "added user with incorrect data")
}

// TestPostService_UpdatePost tests updating a post.
func TestPostService_UpdatePost(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
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
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}
	updatedPost := repository.Post{
		URLHandle: postModel.URLHandle,
		AuthorID:  userModel.ID,
		Title:     postModel.Title,
		Summary:   postModel.Summary,
		Body:      postModel.Body,
	}
	dbErr := fmt.Errorf("error")

	c.mostPostRepository.EXPECT().UpdatePost(updatedPost).Return(postModel, dbErr)

	p, err := c.sut.UpdatePost(updatedPost)

	assert.Equal(t, dbErr, err, "should forward DB error to controller")
	assert.Equal(t, postModel, p, "added post doesn't match the input")
}

// TestPostService_DeletePost tests deleting a post.
func TestPostService_DeletePost(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	urlHandle := "testUrlHandle"
	dbErr := fmt.Errorf("error")

	c.mostPostRepository.EXPECT().DeletePost(urlHandle).Return(dbErr)

	err := c.sut.DeletePost(urlHandle)

	assert.Equal(t, dbErr, err, "should forward DB error to controller")
}

// TestPostService_GetPost tests getting a post from the blog.
func TestPostService_GetPost(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
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
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
		CreatedAt: time.Time{}.Local(),
		UpdatedAt: time.Time{}.Local(),
	}
	post := repository.Post{
		URLHandle: postModel.URLHandle,
		Title:     postModel.Title,
		Author:    userModel,
		Summary:   postModel.Summary,
		Body:      postModel.Body,
		CreatedAt: postModel.CreatedAt,
		UpdatedAt: postModel.UpdatedAt,
	}

	c.mostPostRepository.EXPECT().GetPost(postModel.URLHandle).Return(postModel, nil)

	p, err := c.sut.GetPost(postModel.URLHandle)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, post, p, "post doesn't match the expected output")
}

// TestPostService_GetPost_Unexpected_Error tests handling an unexpected error while getting a post.
func TestPostService_GetPost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	c.mostPostRepository.EXPECT().GetPost("testUrlHandle").Return(repository.Post{}, fmt.Errorf("error"))
	_, err := c.sut.GetPost("testUrlHandle")

	assert.NotNil(t, err, "expected error")
}

// TestPostService_GetPosts tests getting posts from the blog.
func TestPostService_GetPosts(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
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
			Title:     &title,
			Summary:   &summary,
			Body:      &body,
			CreatedAt: time.Time{}.Local(),
			UpdatedAt: time.Time{}.Local(),
		},
	}

	posts := []repository.Post{
		{
			URLHandle: postModels[0].URLHandle,
			Title:     postModels[0].Title,
			Author:    userModel,
			Summary:   postModels[0].Summary,
			CreatedAt: postModels[0].CreatedAt,
			UpdatedAt: postModels[0].UpdatedAt,
			Body:      postModels[0].Body,
		},
	}

	c.mostPostRepository.EXPECT().GetPosts(1, 5).Return(postModels, nil)

	p, err := c.sut.GetPosts()

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, posts, p, "post doesn't match the expected output")
}

// TestPostService_GetPostsPage tests getting a specific page of posts from the blog.
func TestPostService_GetPostsPage(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
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
			Title:     &title,
			Summary:   &summary,
			Body:      &body,
			CreatedAt: time.Time{}.Local(),
			UpdatedAt: time.Time{}.Local(),
		},
	}

	posts := []repository.Post{
		{
			URLHandle: postModels[0].URLHandle,
			Title:     postModels[0].Title,
			Author:    userModel,
			Summary:   postModels[0].Summary,
			CreatedAt: postModels[0].CreatedAt,
			UpdatedAt: postModels[0].UpdatedAt,
			Body:      postModels[0].Body,
		},
	}

	c.mostPostRepository.EXPECT().GetPosts(2, 5).Return(postModels, nil)

	p, err := c.sut.GetPostsPage(2)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, posts, p, "post doesn't match the expected output")
}

// TestPostService_GetPostsPage_Invalid_Page tests getting a specific page of posts from the blog with a negative page number.
func TestPostService_GetPostsPage_Invalid_Page(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	expectedError := errortypes.InvalidPostPageError{Page: -2}
	_, err := c.sut.GetPostsPage(-2)

	assert.Equal(t, expectedError, err, "error doesn't match expected one")
}

// TestPostService_GetPosts_Unexpected_Error tests handling an unexpected error while getting posts
func TestPostService_GetPosts_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostServiceContext(t)

	c.mostPostRepository.EXPECT().GetPosts(1, 5).Return(nil, fmt.Errorf("error"))
	_, err := c.sut.GetPosts()

	assert.NotNil(t, err, "expected error")
}
