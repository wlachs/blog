package controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/blog/api/types"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/controller"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/mocks"
	"github.com/wlachs/blog/internal/repository"
	"github.com/wlachs/blog/internal/test"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"net/url"
	"testing"
)

// postTestContext contains commonly used services, controllers and other objects relevant for testing the PostController.
type postTestContext struct {
	mockPostService *mocks.MockPostService
	sut             controller.PostController
	ctx             *gin.Context
	rec             *httptest.ResponseRecorder
}

// createPostControllerContext creates the context for testing the PostController and reduces code duplication.
func createPostControllerContext(t *testing.T) *postTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockPostService := mocks.NewMockPostService(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), nil, nil, nil)
	sut := controller.CreatePostController(cont, mockPostService)
	ctx, rec := test.CreateControllerContext()

	return &postTestContext{mockPostService, sut, ctx, rec}
}

// TestPostController_AddPost tests adding a new post to the system with valid input params.
func TestPostController_AddPost(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	author := "testAuthor"
	summary := "testSummary"
	body := "testBody"

	input := types.Post{
		Id:      "testUrlHandle",
		Title:   "testTitle",
		Summary: &summary,
		Body:    &body,
	}

	postModel := repository.Post{
		URLHandle: input.Id,
		Title:     &input.Title,
		Summary:   input.Summary,
		Body:      input.Body,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.Set("UserID", author)
	c.ctx.AddParam("PostID", input.Id)
	c.mockPostService.EXPECT().AddPost(postModel, author).Return(postModel, nil)

	c.sut.AddPost(c.ctx)

	var output types.Post
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, input, output, "response body should match")
	assert.Equal(t, 201, c.rec.Code, "incorrect response status")
}

// TestPostController_AddPost_Invalid_Input tests adding a new post to the system with invalid input params.
func TestPostController_AddPost_Invalid_Input(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	c.ctx.Set("user", "testAuthor")

	c.sut.AddPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestPostController_AddPost_Duplicate_Input tests adding a new post to the system with already existing input params.
func TestPostController_AddPost_Duplicate_Input(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
	author := "testAuthor"
	postModel := repository.Post{
		URLHandle: "duplicateUrlHandle",
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}
	input := types.NewPost{
		Title:   postModel.Title,
		Summary: postModel.Summary,
		Body:    postModel.Body,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.Set("UserID", author)
	c.ctx.AddParam("PostID", postModel.URLHandle)
	expectedError := errortypes.DuplicateElementError{}
	c.mockPostService.EXPECT().AddPost(postModel, author).Return(repository.Post{}, expectedError)

	c.sut.AddPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 409, c.rec.Code, "incorrect response status")
}

// TestPostController_AddPost_Unexpected_Error tests handling unexpected errors while adding a new post to the system.
func TestPostController_AddPost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
	userModel := repository.User{
		UserName: "testAuthor",
	}
	postModel := repository.Post{
		URLHandle: "testUrlHandle",
		Author:    userModel,
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}
	inputModel := repository.Post{
		URLHandle: postModel.URLHandle,
		Title:     postModel.Title,
		Summary:   postModel.Summary,
		Body:      postModel.Body,
	}
	input := types.NewPost{
		Body:    postModel.Body,
		Summary: postModel.Summary,
		Title:   postModel.Title,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.Set("UserID", postModel.Author.UserName)
	c.ctx.AddParam("PostID", postModel.URLHandle)
	expectedError := errortypes.UnexpectedPostError{URLHandle: postModel.URLHandle}
	c.mockPostService.EXPECT().AddPost(inputModel, userModel.UserName).Return(repository.Post{}, fmt.Errorf("unexpected internal error"))

	c.sut.AddPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestPostController_UpdatePost tests updating a post with valid input params.
func TestPostController_UpdatePost(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	urlHandle := "testHandle"
	title := "testTitle"
	summary := "testSummary"
	body := "testBody"

	input := types.UpdatedPost{
		Title:   &title,
		Summary: &summary,
		Body:    &body,
	}

	postModel := repository.Post{
		URLHandle: urlHandle,
		Title:     input.Title,
		Summary:   input.Summary,
		Body:      input.Body,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("PostID", urlHandle)
	c.mockPostService.EXPECT().UpdatePost(postModel).Return(postModel, nil)

	c.sut.UpdatePost(c.ctx)

	var output types.UpdatedPost
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, input, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestPostController_UpdatePost_Invalid_Input tests updating a post with invalid input params.
func TestPostController_UpdatePost_Invalid_Input(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	c.sut.UpdatePost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestPostController_UpdatePost_Incorrect_Url_Handle tests updating a non-existing post.
func TestPostController_UpdatePost_Incorrect_Url_Handle(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
	postModel := repository.Post{
		URLHandle: "incorrectUrlHandle",
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}
	input := types.UpdatedPost{
		Title:   postModel.Title,
		Summary: postModel.Summary,
		Body:    postModel.Body,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("PostID", postModel.URLHandle)
	expectedError := errortypes.PostNotFoundError{URLHandle: postModel.URLHandle}
	c.mockPostService.EXPECT().UpdatePost(postModel).Return(repository.Post{}, expectedError)

	c.sut.UpdatePost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 404, c.rec.Code, "incorrect response status")
}

// TestPostController_UpdatePost_Unexpected_Error tests handling unexpected errors while updating a post.
func TestPostController_UpdatePost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
	userModel := repository.User{
		UserName: "testAuthor",
	}
	postModel := repository.Post{
		URLHandle: "testUrlHandle",
		Author:    userModel,
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}
	inputModel := repository.Post{
		URLHandle: postModel.URLHandle,
		Title:     postModel.Title,
		Summary:   postModel.Summary,
		Body:      postModel.Body,
	}
	input := types.UpdatedPost{
		Body:    postModel.Body,
		Summary: postModel.Summary,
		Title:   postModel.Title,
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("PostID", postModel.URLHandle)
	expectedError := errortypes.UnexpectedPostError{URLHandle: postModel.URLHandle}
	c.mockPostService.EXPECT().UpdatePost(inputModel).Return(repository.Post{}, fmt.Errorf("unexpected internal error"))

	c.sut.UpdatePost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestPostController_DeletePost tests deleting a post with valid input params.
func TestPostController_DeletePost(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	urlHandle := "testHandle"

	c.ctx.AddParam("PostID", urlHandle)
	c.mockPostService.EXPECT().DeletePost(urlHandle).Return(nil)

	c.sut.DeletePost(c.ctx)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestPostController_DeletePost_Incorrect_Url_Handle tests deleting a non-existing post.
func TestPostController_DeletePost_Incorrect_Url_Handle(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	urlHandle := "testHandle"

	c.ctx.AddParam("PostID", urlHandle)
	expectedError := errortypes.PostNotFoundError{URLHandle: urlHandle}
	c.mockPostService.EXPECT().DeletePost(urlHandle).Return(expectedError)

	c.sut.DeletePost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 404, c.rec.Code, "incorrect response status")
}

// TestPostController_DeletePost_Unexpected_Error tests handling unexpected errors while deleting a post.
func TestPostController_DeletePost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	urlHandle := "testHandle"

	c.ctx.AddParam("PostID", urlHandle)
	expectedError := errortypes.UnexpectedPostError{URLHandle: urlHandle}
	c.mockPostService.EXPECT().DeletePost(urlHandle).Return(fmt.Errorf("unexpected internal error"))

	c.sut.DeletePost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPost tests retrieving a single post by its URL handle from the blog.
func TestPostController_GetPost(t *testing.T) {
	t.Parallel()
	c := createPostControllerContext(t)

	title := "testTitle"
	summary := "testSummary"
	body := "testBody"
	userModel := repository.User{
		UserName: "testAuthor",
	}
	postModel := repository.Post{
		URLHandle: "testUrlHandle",
		Author:    userModel,
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}
	expectedOutput := types.Post{
		Author:  postModel.Author.UserName,
		Id:      postModel.URLHandle,
		Title:   *postModel.Title,
		Summary: postModel.Summary,
		Body:    postModel.Body,
	}

	c.ctx.AddParam("PostID", postModel.URLHandle)
	c.mockPostService.EXPECT().GetPost(postModel.URLHandle).Return(postModel, nil)

	c.sut.GetPost(c.ctx)

	var output types.Post
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "expected no errors")
	assert.Equal(t, expectedOutput, output, "incorrect output body")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPost_Missing_URL_Handle tests retrieving a single post without URL handle.
func TestPostController_GetPost_Missing_URL_Handle(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	expectedError := errortypes.MissingUrlHandleError{}

	c.sut.GetPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPost_Incorrect_URL_Handle tests retrieving a single post with non-existing URL handle.
func TestPostController_GetPost_Incorrect_URL_Handle(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	urlHandle := "testUrlHandle"
	expectedError := errortypes.PostNotFoundError{URLHandle: urlHandle}

	c.ctx.AddParam("PostID", urlHandle)
	c.mockPostService.EXPECT().GetPost(urlHandle).Return(repository.Post{}, expectedError)

	c.sut.GetPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 404, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPost_Unexpected_Error tests handling an unexpected error while retrieving a single post from the blog.
func TestPostController_GetPost_Unexpected_Error(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	urlHandle := "testUrlHandle"
	expectedError := errortypes.UnexpectedPostError{URLHandle: urlHandle}

	c.ctx.AddParam("PostID", urlHandle)
	c.mockPostService.EXPECT().GetPost(urlHandle).Return(repository.Post{}, fmt.Errorf("unexpected error"))

	c.sut.GetPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPosts tests retrieving a page of posts from the blog.
func TestPostController_GetPosts(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	urlHandle := "testUrlHandle"
	title := "testTitle"
	summary := "testSummary"

	userModel := repository.User{
		UserName: "testAuthor",
	}

	postModels := []repository.Post{
		{
			URLHandle: urlHandle,
			Author:    userModel,
			Title:     &title,
			Summary:   &summary,
		},
	}

	expectedOutput := []types.PostMetadata{
		{
			Id:      urlHandle,
			Title:   title,
			Author:  userModel.UserName,
			Summary: &summary,
		},
	}

	c.mockPostService.EXPECT().GetPosts().Return(postModels, nil)

	c.sut.GetPosts(c.ctx)

	var output []types.PostMetadata
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "expected no errors")
	assert.Equal(t, expectedOutput, output, "incorrect output body")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPostsPage tests retrieving a specific page of posts from the blog.
func TestPostController_GetPostsPage(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	c.ctx.Request.URL, _ = url.Parse("?page=2")

	urlHandle := "testUrlHandle"
	title := "testTitle"
	summary := "testSummary"

	userModel := repository.User{
		UserName: "testAuthor",
	}

	postModels := []repository.Post{
		{
			URLHandle: urlHandle,
			Author:    userModel,
			Title:     &title,
			Summary:   &summary,
		},
	}

	expectedOutput := []types.PostMetadata{
		{
			Id:      urlHandle,
			Title:   title,
			Author:  userModel.UserName,
			Summary: &summary,
		},
	}

	c.mockPostService.EXPECT().GetPostsPage(2).Return(postModels, nil)

	c.sut.GetPosts(c.ctx)

	var output []types.PostMetadata
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "expected no errors")
	assert.Equal(t, expectedOutput, output, "incorrect output body")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPost_Unexpected_Error tests handling an unexpected error while retrieving a single post from the blog.
func TestPostController_GetPosts_Unexpected_Error(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	expectedError := errortypes.UnexpectedPostError{}

	c.mockPostService.EXPECT().GetPosts().Return(nil, fmt.Errorf("unexpected error"))

	c.sut.GetPosts(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}
