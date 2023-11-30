//nolint:paralleltest
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
)

// userTestContext contains objects relevant for testing the UserService.
type userTestContext struct {
	mockUserRepository *mocks.MockUserRepository
	mockJwtUtils       *mocks.MockTokenUtils
	sut                services.UserService
}

// createUserServiceContext creates the context for testing the UserService and reduces code duplication.
func createUserServiceContext(t *testing.T) *userTestContext {
	t.Helper()

	t.Setenv("DEFAULT_USER", "TEST")
	t.Setenv("DEFAULT_PASSWORD", "PW")

	mockCtrl := gomock.NewController(t)
	mockUserRepository := mocks.NewMockUserRepository(mockCtrl)
	mockJwtUtils := mocks.NewMockTokenUtils(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), nil, mockUserRepository, mockJwtUtils)

	mockUserRepository.EXPECT().GetUser("TEST").Return(nil, nil)
	sut := services.CreateUserService(cont)

	return &userTestContext{mockUserRepository, mockJwtUtils, sut}
}

// createUserServiceContext creates the context for testing the UserService and reduces code duplication.
func createUserServiceContextWithoutDefaults(t *testing.T) *userTestContext {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	mockUserRepository := mocks.NewMockUserRepository(mockCtrl)
	mockJwtUtils := mocks.NewMockTokenUtils(mockCtrl)
	cont := container.CreateContainer(logger.CreateLogger(), nil, mockUserRepository, mockJwtUtils)

	sut := services.CreateUserService(cont)

	return &userTestContext{mockUserRepository, mockJwtUtils, sut}
}

// TestUserService_AuthenticateUser tests user authentication.
func TestUserService_AuthenticateUser(t *testing.T) {
	c := createUserServiceContext(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "testAuthor",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	input := types.UserLoginInput{
		UserName: userModel.UserName,
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserName).Return(&userModel, nil)
	c.mockJwtUtils.EXPECT().GenerateJWT(input.UserName).Return("TOKEN", nil)

	token, err := c.sut.AuthenticateUser(&input)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, "TOKEN", token)
}

// TestUserService_AuthenticateUser_Invalid_Password tests user authentication with invalid password.
func TestUserService_AuthenticateUser_Invalid_Password(t *testing.T) {
	c := createUserServiceContext(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "testAuthor",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	input := types.UserLoginInput{
		UserName: userModel.UserName,
		Password: "test",
	}

	expectedError := errortypes.IncorrectUsernameOrPasswordError{}

	c.mockUserRepository.EXPECT().GetUser(input.UserName).Return(&userModel, nil)

	token, err := c.sut.AuthenticateUser(&input)

	assert.Equal(t, token, "", "no token should be generated")
	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_CheckUserPassword tests checking the user's password upon login.
func TestUserService_CheckUserPassword(t *testing.T) {
	c := createUserServiceContext(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "testAuthor",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	input := types.UserLoginInput{
		UserName: userModel.UserName,
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserName).Return(&userModel, nil)

	success := c.sut.CheckUserPassword(&input)

	assert.True(t, success, "password should match the one stored in the database")
}

// TestUserService_CheckUserPassword_Invalid_User tests checking the user's password upon login with incorrect username.
func TestUserService_CheckUserPassword_Invalid_User(t *testing.T) {
	c := createUserServiceContext(t)

	input := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserName).Return(nil, fmt.Errorf("internal error"))

	success := c.sut.CheckUserPassword(&input)

	assert.False(t, success, "username should not match the one stored in the database")
}

// TestUserService_CheckUserPassword_Invalid_Password tests checking the user's password upon login with false credentials.
func TestUserService_CheckUserPassword_Invalid_Password(t *testing.T) {
	c := createUserServiceContext(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "testAuthor",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	input := types.UserLoginInput{
		UserName: userModel.UserName,
		Password: "test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserName).Return(&userModel, nil)

	success := c.sut.CheckUserPassword(&input)

	assert.False(t, success, "password should not match the one stored in the database")
}

// TestUserService_GetUser tests getting a single user of the blog.
func TestUserService_GetUser(t *testing.T) {
	c := createUserServiceContext(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "testAuthor",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts: []repository.Post{{
			URLHandle: "handle",
		}},
	}

	expectedUser := types.User{
		UserName:     userModel.UserName,
		PasswordHash: userModel.PasswordHash,
		Posts:        []string{userModel.Posts[0].URLHandle},
	}

	c.mockUserRepository.EXPECT().GetUser(userModel.UserName).Return(&userModel, nil)

	user, err := c.sut.GetUser(userModel.UserName)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, expectedUser, user, "response doesn't match expected user data")
}

// TestUserService_GetUser_Unexpected_Error tests handling an unexpected error while getting a single user of the blog.
func TestUserService_GetUser_Unexpected_Error(t *testing.T) {
	c := createUserServiceContext(t)

	c.mockUserRepository.EXPECT().GetUser(gomock.Any()).Return(nil, fmt.Errorf("internal error"))

	user, err := c.sut.GetUser("")

	assert.NotNil(t, err, "expected to receive and error")
	assert.Equal(t, types.User{}, user, "response doesn't match expected user data")
}

// TestUserService_GetUsers tests getting every user of the blog.
func TestUserService_GetUsers(t *testing.T) {
	c := createUserServiceContext(t)

	userModels := []repository.User{
		{
			ID:           0,
			UserName:     "testAuthor1",
			PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
			Posts: []repository.Post{{
				URLHandle: "handle1",
			}},
		},
		{
			ID:           1,
			UserName:     "testAuthor2",
			PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
			Posts: []repository.Post{{
				URLHandle: "handle2",
			}},
		},
	}

	expectedUsers := []types.User{
		{
			UserName:     userModels[0].UserName,
			PasswordHash: userModels[0].PasswordHash,
			Posts:        []string{userModels[0].Posts[0].URLHandle},
		},
		{
			UserName:     userModels[1].UserName,
			PasswordHash: userModels[1].PasswordHash,
			Posts:        []string{userModels[1].Posts[0].URLHandle},
		},
	}

	c.mockUserRepository.EXPECT().GetUsers().Return(userModels, nil)

	users, err := c.sut.GetUsers()

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, expectedUsers, users, "response doesn't match expected user data")
}

// TestUserService_GetUsers_Unexpected_Error tests handling an unexpected error while getting every user of the blog.
func TestUserService_GetUsers_Unexpected_Error(t *testing.T) {
	c := createUserServiceContext(t)

	c.mockUserRepository.EXPECT().GetUsers().Return(nil, fmt.Errorf("internal error"))

	users, err := c.sut.GetUsers()

	assert.NotNil(t, err, "expected to receive an error")
	assert.Equal(t, []types.User{}, users, "response doesn't match expected user data")
}

// TestUserService_RegisterFirstUser tests registering the first user.
func TestUserService_RegisterFirstUser(t *testing.T) {
	c := createUserServiceContextWithoutDefaults(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "TEST",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	t.Setenv("DEFAULT_USER", userModel.UserName)
	t.Setenv("DEFAULT_PASSWORD", "Test")

	c.mockUserRepository.EXPECT().GetUser(userModel.UserName).Return(nil, fmt.Errorf("internal error"))
	c.mockUserRepository.EXPECT().AddUser(gomock.Any()).Return(&userModel, nil)

	err := c.sut.RegisterFirstUser()

	assert.Nil(t, err, "expected to complete without error")
}

// TestUserService_RegisterFirstUser_Invalid_Defaults tests registering the first user with missing credentials.
func TestUserService_RegisterFirstUser_Invalid_Defaults(t *testing.T) {
	t.Parallel()
	c := createUserServiceContextWithoutDefaults(t)

	expectedError := errortypes.MissingDefaultUsernameOrPasswordError{}

	err := c.sut.RegisterFirstUser()

	assert.Equal(t, expectedError, err, "error type doesn't match")
}

// TestUserService_RegisterUser tests adding a new user to the system.
func TestUserService_RegisterUser(t *testing.T) {
	c := createUserServiceContext(t)

	userModel := repository.User{
		ID:           0,
		UserName:     "testAuthor",
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	input := types.UserLoginInput{
		UserName: userModel.UserName,
		Password: "Test",
	}

	expectedUser := types.User{
		UserName:     userModel.UserName,
		PasswordHash: userModel.PasswordHash,
		Posts:        []string{},
	}

	c.mockUserRepository.EXPECT().AddUser(gomock.Any()).Return(&userModel, nil)

	user, err := c.sut.RegisterUser(&input)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, expectedUser, user, "response doesn't match expected user data")
}

// TestUserService_RegisterUser_Invalid_Password tests adding a new user to the system with a password too long.
func TestUserService_RegisterUser_Invalid_Password(t *testing.T) {
	c := createUserServiceContext(t)

	input := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "1234567890123456789012345678901234567890123456789012345678901234567890123",
	}

	expectedError := errortypes.PasswordHashingError{}

	_, err := c.sut.RegisterUser(&input)

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_UpdateUser tests updating an existing user.
func TestUserService_UpdateUser(t *testing.T) {
	c := createUserServiceContext(t)

	oldUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test",
	}

	newUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test1",
	}

	oldUserModel := repository.User{
		ID:           0,
		UserName:     oldUser.UserName,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	newUserModel := repository.User{
		ID:           0,
		UserName:     oldUser.UserName,
		PasswordHash: "$2y$10$o7ZqUckxyaZAS31yFfBhTutbo3cWUQkvsdnVikvhrn69.c5kG0/TS",
		Posts:        []repository.Post{},
	}

	expectedNewUser := types.User{
		UserName:     oldUser.UserName,
		PasswordHash: newUserModel.PasswordHash,
		Posts:        []string{},
	}

	c.mockUserRepository.EXPECT().GetUser(oldUser.UserName).Return(&oldUserModel, nil)
	c.mockUserRepository.EXPECT().UpdateUser(gomock.Any()).Return(&newUserModel, nil)

	user, err := c.sut.UpdateUser(&oldUser, &newUser)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, expectedNewUser, user, "response doesn't match expected user data")
}

// TestUserService_UpdateUser_Invalid_Old_Password tests updating an existing user with an incorrect password.
func TestUserService_UpdateUser_Invalid_Old_Password(t *testing.T) {
	c := createUserServiceContext(t)

	oldUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test1",
	}

	oldUserModel := repository.User{
		ID:           0,
		UserName:     oldUser.UserName,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	expectedError := errortypes.IncorrectUsernameOrPasswordError{}

	c.mockUserRepository.EXPECT().GetUser(oldUser.UserName).Return(&oldUserModel, nil)

	_, err := c.sut.UpdateUser(&oldUser, &types.UserLoginInput{})

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_UpdateUser_Invalid_New_Password tests updating an existing user with a password too long.
func TestUserService_UpdateUser_Invalid_New_Password(t *testing.T) {
	c := createUserServiceContext(t)

	oldUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test",
	}

	newUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "1234567890123456789012345678901234567890123456789012345678901234567890123",
	}

	oldUserModel := repository.User{
		ID:           0,
		UserName:     oldUser.UserName,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	expectedError := errortypes.PasswordHashingError{}

	c.mockUserRepository.EXPECT().GetUser(oldUser.UserName).Return(&oldUserModel, nil)

	_, err := c.sut.UpdateUser(&oldUser, &newUser)

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_UpdateUser_Unexpected_Error tests handling errors while updating an existing user.
func TestUserService_UpdateUser_Unexpected_Error(t *testing.T) {
	c := createUserServiceContext(t)

	oldUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test",
	}

	newUser := types.UserLoginInput{
		UserName: "testAuthor",
		Password: "Test1",
	}

	oldUserModel := repository.User{
		ID:           0,
		UserName:     oldUser.UserName,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	c.mockUserRepository.EXPECT().GetUser(oldUser.UserName).Return(&oldUserModel, nil)
	c.mockUserRepository.EXPECT().UpdateUser(gomock.Any()).Return(nil, fmt.Errorf("internal error"))

	_, err := c.sut.UpdateUser(&oldUser, &newUser)

	assert.NotNil(t, err, "expected to receive an error")
}
