//nolint:paralleltest
package services_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/blog/api/types"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/mocks"
	"github.com/wlachs/blog/internal/repository"
	"github.com/wlachs/blog/internal/services"
	"go.uber.org/mock/gomock"
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

	mockUserRepository.EXPECT().GetUser("TEST").Return(repository.User{}, nil)
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

	input := types.DoLoginJSONBody{
		UserID:   userModel.UserName,
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserID).Return(userModel, nil)
	c.mockJwtUtils.EXPECT().GenerateJWT(input.UserID).Return("TOKEN", nil)

	token, err := c.sut.AuthenticateUser(input.UserID, input.Password)

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

	input := types.DoLoginJSONBody{
		UserID:   userModel.UserName,
		Password: "test",
	}

	expectedError := errortypes.IncorrectUsernameOrPasswordError{}

	c.mockUserRepository.EXPECT().GetUser(input.UserID).Return(userModel, nil)

	token, err := c.sut.AuthenticateUser(input.UserID, input.Password)

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

	input := types.DoLoginJSONBody{
		UserID:   userModel.UserName,
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserID).Return(userModel, nil)

	success := c.sut.CheckUserPassword(input.UserID, input.Password)

	assert.True(t, success, "password should match the one stored in the database")
}

// TestUserService_CheckUserPassword_Invalid_User tests checking the user's password upon login with incorrect username.
func TestUserService_CheckUserPassword_Invalid_User(t *testing.T) {
	c := createUserServiceContext(t)

	input := types.DoLoginJSONBody{
		UserID:   "testAuthor",
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserID).Return(repository.User{}, fmt.Errorf("internal error"))

	success := c.sut.CheckUserPassword(input.UserID, input.Password)

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

	input := types.DoLoginJSONBody{
		UserID:   userModel.UserName,
		Password: "test",
	}

	c.mockUserRepository.EXPECT().GetUser(input.UserID).Return(userModel, nil)

	success := c.sut.CheckUserPassword(input.UserID, input.Password)

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

	c.mockUserRepository.EXPECT().GetUser(userModel.UserName).Return(userModel, nil)

	user, err := c.sut.GetUser(userModel.UserName)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, userModel, user, "response doesn't match expected user data")
}

// TestUserService_GetUser_Unexpected_Error tests handling an unexpected error while getting a single user of the blog.
func TestUserService_GetUser_Unexpected_Error(t *testing.T) {
	c := createUserServiceContext(t)

	c.mockUserRepository.EXPECT().GetUser(gomock.Any()).Return(repository.User{}, fmt.Errorf("internal error"))

	user, err := c.sut.GetUser("")

	assert.NotNil(t, err, "expected to receive and error")
	assert.Equal(t, repository.User{}, user, "response doesn't match expected user data")
}

// TestUserService_GetUsers tests getting the first page of users of the blog.
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

	c.mockUserRepository.EXPECT().GetUsers(1, 5).Return(userModels, 1, nil)

	users, _, err := c.sut.GetUsers()

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, userModels, users, "response doesn't match expected user data")
}

// TestUserService_GetUsers_Unexpected_Error tests handling an unexpected error while getting every user of the blog.
func TestUserService_GetUsers_Unexpected_Error(t *testing.T) {
	c := createUserServiceContext(t)

	c.mockUserRepository.EXPECT().GetUsers(1, 5).Return([]repository.User{}, -1, fmt.Errorf("internal error"))

	users, _, err := c.sut.GetUsers()

	assert.NotNil(t, err, "expected to receive an error")
	assert.Equal(t, []repository.User{}, users, "response doesn't match expected user data")
}

// TestUserService_GetUsersPage tests getting a specific page of users of the blog.
func TestUserService_GetUsersPage(t *testing.T) {
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

	c.mockUserRepository.EXPECT().GetUsers(2, 5).Return(userModels, 2, nil)

	users, _, err := c.sut.GetUsersPage(2)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, userModels, users, "response doesn't match expected user data")
}

// TestUserService_GetUsersPage_Invalid_Page tests getting a specific page of users of the blog with a negative page number.
func TestUserService_GetUsersPage_Invalid_Page(t *testing.T) {
	c := createUserServiceContext(t)

	expectedError := errortypes.InvalidUserPageError{Page: -2}
	_, _, err := c.sut.GetUsersPage(-2)

	assert.Equal(t, expectedError, err, "error doesn't match expected one")
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

	c.mockUserRepository.EXPECT().GetUser(userModel.UserName).Return(repository.User{}, fmt.Errorf("internal error"))
	c.mockUserRepository.EXPECT().AddUser(gomock.Any()).Return(userModel, nil)

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

	input := types.DoLoginJSONBody{
		UserID:   userModel.UserName,
		Password: "Test",
	}

	c.mockUserRepository.EXPECT().AddUser(gomock.Any()).Return(userModel, nil)

	user, err := c.sut.RegisterUser(input.UserID, input.Password)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, userModel, user, "response doesn't match expected user data")
}

// TestUserService_RegisterUser_Invalid_Password tests adding a new user to the system with a password too long.
func TestUserService_RegisterUser_Invalid_Password(t *testing.T) {
	c := createUserServiceContext(t)

	input := types.DoLoginJSONBody{
		UserID:   "testAuthor",
		Password: "1234567890123456789012345678901234567890123456789012345678901234567890123",
	}

	expectedError := errortypes.PasswordHashingError{}

	_, err := c.sut.RegisterUser(input.UserID, input.Password)

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_RegisterUser_Missing_Password tests adding a new user to the system without a password.
func TestUserService_RegisterUser_Missing_Password(t *testing.T) {
	c := createUserServiceContext(t)

	input := types.DoLoginJSONBody{
		UserID:   "testAuthor",
		Password: "",
	}

	expectedError := errortypes.MissingPasswordError{}

	_, err := c.sut.RegisterUser(input.UserID, input.Password)

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_UpdateUser tests updating an existing user.
func TestUserService_UpdateUser(t *testing.T) {
	c := createUserServiceContext(t)

	userID := "testAuthor"
	oldPassword := "Test"
	newPassword := "Test1"
	oldUserModel := repository.User{
		ID:           0,
		UserName:     userID,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	newUserModel := repository.User{
		ID:           0,
		UserName:     userID,
		PasswordHash: "$2y$10$o7ZqUckxyaZAS31yFfBhTutbo3cWUQkvsdnVikvhrn69.c5kG0/TS",
		Posts:        []repository.Post{},
	}

	c.mockUserRepository.EXPECT().GetUser(userID).Return(oldUserModel, nil)
	c.mockUserRepository.EXPECT().UpdateUser(gomock.Any()).Return(newUserModel, nil)

	user, err := c.sut.UpdateUser(userID, oldPassword, newPassword)

	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, newUserModel, user, "response doesn't match expected user data")
}

// TestUserService_UpdateUser_Invalid_Old_Password tests updating an existing user with an incorrect password.
func TestUserService_UpdateUser_Invalid_Old_Password(t *testing.T) {
	c := createUserServiceContext(t)

	userID := "testAuthor"
	oldPassword := "wrong password"
	newPassword := "something new"

	oldUserModel := repository.User{
		ID:           0,
		UserName:     userID,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	expectedError := errortypes.IncorrectUsernameOrPasswordError{}

	c.mockUserRepository.EXPECT().GetUser(userID).Return(oldUserModel, nil)

	_, err := c.sut.UpdateUser(userID, oldPassword, newPassword)

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_UpdateUser_Invalid_New_Password tests updating an existing user with a password too long.
func TestUserService_UpdateUser_Invalid_New_Password(t *testing.T) {
	c := createUserServiceContext(t)

	userID := "testAuthor"
	oldPassword := "Test"
	newPassword := "1234567890123456789012345678901234567890123456789012345678901234567890123"

	oldUserModel := repository.User{
		ID:           0,
		UserName:     userID,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	expectedError := errortypes.PasswordHashingError{}

	c.mockUserRepository.EXPECT().GetUser(userID).Return(oldUserModel, nil)

	_, err := c.sut.UpdateUser(userID, oldPassword, newPassword)

	assert.Equal(t, expectedError, err, "incorrect error type")
}

// TestUserService_UpdateUser_Unexpected_Error tests handling errors while updating an existing user.
func TestUserService_UpdateUser_Unexpected_Error(t *testing.T) {
	c := createUserServiceContext(t)

	userID := "testAuthor"
	oldPassword := "Test"
	newPassword := "Test1"

	oldUserModel := repository.User{
		ID:           0,
		UserName:     userID,
		PasswordHash: "$2y$10$Hb7smnjLlPtN.VMyNi5dYuMaCmEgCbus/Tapxf2u5jhxkKE1Pr50.",
		Posts:        []repository.Post{},
	}

	c.mockUserRepository.EXPECT().GetUser(userID).Return(oldUserModel, nil)
	c.mockUserRepository.EXPECT().UpdateUser(gomock.Any()).Return(repository.User{}, fmt.Errorf("internal error"))

	_, err := c.sut.UpdateUser(userID, oldPassword, newPassword)

	assert.NotNil(t, err, "expected to receive an error")
}

// TestUserService_DeleteUser tests deleting an existing user.
func TestUserService_DeleteUser(t *testing.T) {
	c := createUserServiceContext(t)

	userID := "testAuthor"
	dbErr := fmt.Errorf("unexpected error")

	c.mockUserRepository.EXPECT().DeleteUser(gomock.Any()).Return(dbErr)

	err := c.sut.DeleteUser(userID)

	assert.Equal(t, dbErr, err, "should forward DB error to controller")
}
