package controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/test"
	"github.com/wlchs/blog/internal/types"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
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

	input := types.Post{
		URLHandle: "testUrlHandle",
		Title:     "testTitle",
		Author:    "testAuthor",
		Summary:   "testSummary",
		Body:      "testBody",
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.Set("user", input.Author)
	c.mockPostService.EXPECT().AddPost(&input).Return(input, nil)

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

	input := types.Post{
		URLHandle: "duplicateUrlHandle",
		Title:     "testTitle",
		Author:    "testAuthor",
		Summary:   "testSummary",
		Body:      "testBody",
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.Set("user", input.Author)
	expectedError := errortypes.DuplicateElementError{}
	c.mockPostService.EXPECT().AddPost(&input).Return(types.Post{}, expectedError)

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

	input := types.Post{
		URLHandle: "testUrlHandle",
		Title:     "testTitle",
		Author:    "testAuthor",
		Summary:   "testSummary",
		Body:      "testBody",
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.Set("user", input.Author)
	expectedError := errortypes.UnexpectedPostError{Post: input}
	c.mockPostService.EXPECT().AddPost(&input).Return(types.Post{}, fmt.Errorf("unexpected internal error"))

	c.sut.AddPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPost tests retrieving a single post by its URL handle from the blog.
func TestPostController_GetPost(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	expectedOutput := types.Post{
		URLHandle: "testUrlHandle",
		Title:     "testTitle",
		Author:    "testAuthor",
		Summary:   "testSummary",
		Body:      "testBody",
	}

	c.ctx.AddParam("id", expectedOutput.URLHandle)
	c.mockPostService.EXPECT().GetPost(expectedOutput.URLHandle).Return(expectedOutput, nil)

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
	expectedError := errortypes.PostNotFoundError{Post: types.Post{URLHandle: urlHandle}}

	c.ctx.AddParam("id", urlHandle)
	c.mockPostService.EXPECT().GetPost(urlHandle).Return(types.Post{}, expectedError)

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
	expectedError := errortypes.UnexpectedPostError{Post: types.Post{URLHandle: urlHandle}}

	c.ctx.AddParam("id", urlHandle)
	c.mockPostService.EXPECT().GetPost(urlHandle).Return(types.Post{}, fmt.Errorf("unexpected error"))

	c.sut.GetPost(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestPostController_GetPosts tests retrieving a every post from the blog.
func TestPostController_GetPosts(t *testing.T) {
	t.Parallel()

	c := createPostControllerContext(t)
	expectedOutput := []types.Post{
		{
			URLHandle: "testUrlHandle",
			Title:     "testTitle",
			Author:    "testAuthor",
			Summary:   "testSummary",
			Body:      "testBody",
		},
	}

	c.mockPostService.EXPECT().GetPosts().Return(expectedOutput, nil)

	c.sut.GetPosts(c.ctx)

	var output []types.Post
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
