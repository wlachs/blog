package controller_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/controller"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/mocks"
	"github.com/wlchs/blog/internal/test"
	"github.com/wlchs/blog/internal/types"
	"testing"
)

// TestAuthController_Login tests the login method on the AuthController with valid data.
func TestAuthController_Login(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := mocks.NewMockContainer(mockCtrl)
	mockUserService := mocks.NewMockUserService(mockCtrl)

	sut := controller.CreateAuthController(mockContainer, mockUserService)
	ctx, rec := test.CreateControllerContext()

	input := types.UserLoginInput{
		UserName: "TestUser",
		Password: "TestPW1234$",
	}

	test.MockJsonPost(ctx, map[string]interface{}{
		"userName": input.UserName,
		"password": input.Password,
	})

	mockUserService.EXPECT().AuthenticateUser(&input).Return("token", nil)

	sut.Login(ctx)
	assert.Nil(t, ctx.Errors, "should complete without errors")
	assert.Equal(t, 200, rec.Code)
}

// TestAuthController_Login_Incorrect_Password tests the login method on the AuthController with valid data but incorrect password.
func TestAuthController_Login_Incorrect_Password(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := mocks.NewMockContainer(mockCtrl)
	mockUserService := mocks.NewMockUserService(mockCtrl)

	sut := controller.CreateAuthController(mockContainer, mockUserService)
	ctx, rec := test.CreateControllerContext()

	input := types.UserLoginInput{
		UserName: "TestUser",
		Password: "TestPW1234$",
	}

	test.MockJsonPost(ctx, map[string]interface{}{
		"userName": input.UserName,
		"password": input.Password,
	})

	expectedError := errortypes.IncorrectUsernameOrPasswordError{}
	mockUserService.EXPECT().AuthenticateUser(&input).Return("", expectedError)

	sut.Login(ctx)

	errors := ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors), "expected one error")
	assert.Equal(t, expectedError.Error(), errors[0])
	assert.Equal(t, 401, rec.Code)
}

// TestAuthController_Login_Invalid_Input tests the login method on the AuthController with invalid data.
func TestAuthController_Login_Invalid_Input(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockContainer := mocks.NewMockContainer(mockCtrl)
	mockUserService := mocks.NewMockUserService(mockCtrl)

	sut := controller.CreateAuthController(mockContainer, mockUserService)
	ctx, rec := test.CreateControllerContext()

	sut.Login(ctx)

	errors := ctx.Errors.Errors()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, 400, rec.Code)
}
