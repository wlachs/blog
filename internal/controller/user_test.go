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
	"testing"
)

// userTestContext contains commonly used services, controllers and other objects relevant for testing the UserController.
type userTestContext struct {
	mockUserService *mocks.MockUserService
	sut             controller.UserController
	ctx             *gin.Context
	rec             *httptest.ResponseRecorder
}

// createUserControllerContext creates the context for testing the UserController and reduces code duplication.
func createUserControllerContext(t *testing.T) *userTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockUserService := mocks.NewMockUserService(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), nil, nil, nil)
	sut := controller.CreateUserController(cont, mockUserService)
	ctx, rec := test.CreateControllerContext()

	return &userTestContext{mockUserService, sut, ctx, rec}
}

// TestUserController_GetUser tests retrieving a user from the blog.
func TestUserController_GetUser(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userModel := repository.User{
		UserName: "testAuthor",
	}
	expectedOutput := types.User{
		UserID: userModel.UserName,
	}

	c.ctx.AddParam("UserID", expectedOutput.UserID)
	c.mockUserService.EXPECT().GetUser(expectedOutput.UserID).Return(userModel, nil)

	c.sut.GetUser(c.ctx)

	var output types.User
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, expectedOutput, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUser_Missing_Username tests retrieving a user from the blog without username.
func TestUserController_GetUser_Missing_Username(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	expectedError := errortypes.MissingUsernameError{}

	c.sut.GetUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUser_Incorrect_Username tests retrieving a non-existing user from the blog.
func TestUserController_GetUser_Incorrect_Username(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	expectedError := errortypes.UserNotFoundError{UserName: userName}
	c.ctx.AddParam("UserID", userName)
	c.mockUserService.EXPECT().GetUser(userName).Return(repository.User{}, expectedError)

	c.sut.GetUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 404, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUser_Unexpected_Error tests handling an unexpected error while retrieving a user from the blog.
func TestUserController_GetUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	expectedError := errortypes.UnexpectedUserError{UserName: userName}
	c.ctx.AddParam("UserID", userName)
	c.mockUserService.EXPECT().GetUser(userName).Return(repository.User{}, fmt.Errorf("unexpected error"))

	c.sut.GetUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUsers tests retrieving every user from the blog.
func TestUserController_GetUsers(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userModels := []repository.User{
		{
			UserName: "testAuthor",
			Posts: []repository.Post{
				{
					Author:    repository.User{UserName: "testAuthor"},
					URLHandle: "urlHandle1",
				},
				{
					Author:    repository.User{UserName: "testAuthor"},
					URLHandle: "urlHandle2",
				},
			},
		},
	}
	expectedOutput := []types.User{
		{
			UserID: userModels[0].UserName,
			Posts: &[]types.PostMetadata{
				{
					Id:     userModels[0].Posts[0].URLHandle,
					Author: userModels[0].UserName,
				},
				{
					Id:     userModels[0].Posts[1].URLHandle,
					Author: userModels[0].UserName,
				},
			},
		},
	}

	c.mockUserService.EXPECT().GetUsers().Return(userModels, nil)

	c.sut.GetUsers(c.ctx)

	var output []types.User
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, expectedOutput, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestUserController_GetUsers_Unexpected_Error tests handling an unexpected error while retrieving users from the blog.
func TestUserController_GetUsers_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	expectedError := errortypes.UnexpectedUserError{}
	c.mockUserService.EXPECT().GetUsers().Return(nil, fmt.Errorf("unexpected error"))

	c.sut.GetUsers(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser tests updating a user's password.
func TestUserController_UpdateUser(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	input := types.UpdateUserJSONBody{
		OldPassword: "oldPW",
		NewPassword: "newPW",
	}
	userModel := repository.User{
		UserName: "testAuthor",
		Posts: []repository.Post{
			{
				Author:    repository.User{UserName: "testAuthor"},
				URLHandle: "urlHandle1",
			},
			{
				Author:    repository.User{UserName: "testAuthor"},
				URLHandle: "urlHandle2",
			},
		},
	}
	expectedOutput := types.User{
		UserID: userModel.UserName,
		Posts: &[]types.PostMetadata{
			{
				Id:     userModel.Posts[0].URLHandle,
				Author: userModel.UserName,
			},
			{
				Id:     userModel.Posts[1].URLHandle,
				Author: userModel.UserName,
			},
		},
	}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("UserID", expectedOutput.UserID)
	c.mockUserService.EXPECT().UpdateUser(expectedOutput.UserID, input.OldPassword, input.NewPassword).Return(userModel, nil)
	c.sut.UpdateUser(c.ctx)

	var output types.User
	_ = json.Unmarshal(c.rec.Body.Bytes(), &output)

	assert.Nil(t, c.ctx.Errors, "should complete without error")
	assert.Equal(t, expectedOutput, output, "response body should match")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser_Invalid_Input tests updating a user's password with invalid input.
func TestUserController_UpdateUser_Invalid_Input(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	c.ctx.AddParam("UserID", userName)

	c.sut.UpdateUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser_Incorrect_Username tests updating a user's password with incorrect old password.
func TestUserController_UpdateUser_Incorrect_Username(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	input := types.UpdateUserJSONBody{
		OldPassword: "oldPW",
		NewPassword: "newPW",
	}
	expectedError := errortypes.IncorrectUsernameOrPasswordError{}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("UserID", userName)
	c.mockUserService.EXPECT().UpdateUser(userName, input.OldPassword, input.NewPassword).Return(repository.User{}, expectedError)
	c.sut.UpdateUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 401, c.rec.Code, "incorrect response status")
}

// TestUserController_UpdateUser_Unexpected_Error tests handling an unexpected error while updating a user's password.
func TestUserController_UpdateUser_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createUserControllerContext(t)

	userName := "testAuthor"
	input := types.UpdateUserJSONBody{
		OldPassword: "oldPW",
		NewPassword: "newPW",
	}
	expectedError := errortypes.UnexpectedUserError{UserName: userName}

	test.MockJsonPost(c.ctx, input)

	c.ctx.AddParam("UserID", userName)
	c.mockUserService.EXPECT().UpdateUser(userName, input.OldPassword, input.NewPassword).Return(repository.User{}, fmt.Errorf("unexpected error"))
	c.sut.UpdateUser(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 500, c.rec.Code, "incorrect response status")
}
