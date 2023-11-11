package controller_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/test"
	"github.com/wlchs/blog/internal/types"
	"net/http/httptest"
	"testing"
)

// authTestContext contains commonly used services, controllers and other objects relevant for testing the AuthController.
type authTestContext struct {
	mockUserService *mocks.MockUserService
	sut             controller.AuthController
	ctx             *gin.Context
	rec             *httptest.ResponseRecorder
}

// createAuthControllerContext creates the context for testing the AuthController and reduces code duplication.
func createAuthControllerContext(t *testing.T) *authTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockContainer := mocks.NewMockContainer(mockCtrl)
	mockUserService := mocks.NewMockUserService(mockCtrl)
	sut := controller.CreateAuthController(mockContainer, mockUserService)
	ctx, rec := test.CreateControllerContext()

	return &authTestContext{mockUserService, sut, ctx, rec}
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
	assert.Equal(t, 200, c.rec.Code)
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
	assert.Equal(t, expectedError.Error(), errors[0])
	assert.Equal(t, 401, c.rec.Code)
}

// TestAuthController_Login_Invalid_Input tests the login method on the AuthController with invalid data.
func TestAuthController_Login_Invalid_Input(t *testing.T) {
	t.Parallel()
	c := createAuthControllerContext(t)

	c.sut.Login(c.ctx)

	errors := c.ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, 400, c.rec.Code)
}
