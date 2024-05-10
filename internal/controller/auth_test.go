package controller_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/controller"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/mocks"
	"github.com/wlachs/blog/internal/test"
	"github.com/wlachs/blog/internal/types"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
)

// authTestContext contains commonly used services, controllers and other objects relevant for testing the AuthController.
type authTestContext struct {
	mockUserService *mocks.MockUserService
	mockJwtUtils    *mocks.MockTokenUtils
	sut             controller.AuthController
	ctx             *gin.Context
	rec             *httptest.ResponseRecorder
}

// createAuthControllerContext creates the context for testing the AuthController and reduces code duplication.
func createAuthControllerContext(t *testing.T) *authTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockJwtUtils := mocks.NewMockTokenUtils(mockCtrl)
	mockUserService := mocks.NewMockUserService(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), nil, nil, mockJwtUtils)
	sut := controller.CreateAuthController(cont, mockUserService)
	ctx, rec := test.CreateControllerContext()

	return &authTestContext{mockUserService, mockJwtUtils, sut, ctx, rec}
}

// TestAuthController_Login tests the login method on the AuthController with valid data.
func TestAuthController_Login(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	input := types.UserLoginInput{
		UserName: "TestUser",
		Password: "TestPW1234$",
	}

	test.MockJsonPost(c.ctx, map[string]interface{}{
		"userName": input.UserName,
		"password": input.Password,
	})

	c.mockUserService.EXPECT().AuthenticateUser(&input).Return("token", nil)

	c.sut.Login(c.ctx)
	assert.Nil(t, c.ctx.Errors, "should complete without errors")
	assert.Equal(t, "token", c.rec.Header().Get("X-Auth-Token"))
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestAuthController_Login_Incorrect_Password tests the login method on the AuthController with valid data but incorrect password.
func TestAuthController_Login_Incorrect_Password(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	input := types.UserLoginInput{
		UserName: "TestUser",
		Password: "TestPW1234$",
	}

	test.MockJsonPost(c.ctx, map[string]interface{}{
		"userName": input.UserName,
		"password": input.Password,
	})

	expectedError := errortypes.IncorrectUsernameOrPasswordError{}
	c.mockUserService.EXPECT().AuthenticateUser(&input).Return("", expectedError)

	c.sut.Login(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected one error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 401, c.rec.Code, "incorrect response status")
}

// TestAuthController_Login_Invalid_Input tests the login method on the AuthController with invalid data.
func TestAuthController_Login_Invalid_Input(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	c.sut.Login(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, 400, c.rec.Code, "incorrect response status")
}

// TestAuthController_Protect tests the protect middleware of the AuthController with valid input.
func TestAuthController_Protect(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	c.ctx.Request.Header.Add("X-Auth-Token", "token")
	c.mockJwtUtils.EXPECT().ParseJWT("token").Return("test user", nil)

	c.sut.Protect(c.ctx)

	assert.Nil(t, c.ctx.Errors, "expected no errors")
	assert.Equal(t, "test user", c.ctx.GetString("user"), "incorrect user")
	assert.Equal(t, 200, c.rec.Code, "incorrect response status")
}

// TestAuthController_Protect_Token_Missing tests the protect middleware of the AuthController with missing token.
func TestAuthController_Protect_Token_Missing(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	expectedError := errortypes.MissingAuthTokenError{}

	c.sut.Protect(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 401, c.rec.Code, "incorrect response status")
}

// TestAuthController_Protect_Token_Invalid tests the protect middleware of the AuthController with an invalid token.
func TestAuthController_Protect_Token_Invalid(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	expectedError := errortypes.InvalidAuthTokenError{}
	c.ctx.Request.Header.Add("X-Auth-Token", "token")
	c.mockJwtUtils.EXPECT().ParseJWT("token").Return("", fmt.Errorf("internal error"))

	c.sut.Protect(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected exactly 1 error")
	assert.Equal(t, expectedError.Error(), errors[0], "incorrect error type")
	assert.Equal(t, 401, c.rec.Code, "incorrect response status")
}
