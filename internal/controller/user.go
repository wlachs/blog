package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/services"
	"github.com/wlachs/blog/internal/types"
	"net/http"
)

// UserController interface defining user-related middleware methods to handler HTTP requests.
type UserController interface {
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
}

// userController is a concrete implementation of the UserController interface.
type userController struct {
	cont        container.Container
	userService services.UserService
}

// CreateUserController instantiates a user controller user the application container.
func CreateUserController(cont container.Container, userService services.UserService) UserController {
	return &userController{cont, userService}
}

// GetUser middleware. Top level handler of /user/:userName GET requests.
func (u userController) GetUser(c *gin.Context) {
	userService := u.userService
	userName, found := c.Params.Get("userName")

	if !found {
		_ = c.AbortWithError(http.StatusBadRequest, errortypes.MissingUsernameError{})
		return
	}

	user, err := userService.GetUser(userName)
	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusOK, user)

	case errortypes.UserNotFoundError:
		_ = c.AbortWithError(http.StatusNotFound, err)

	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedUserError{User: types.User{UserName: userName}})
	}
}

// GetUsers middleware. Top level handler of /users GET requests.
func (u userController) GetUsers(c *gin.Context) {
	userService := u.userService
	users, err := userService.GetUsers()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedUserError{})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// UpdateUser middleware. Top level handler of /users/:userName PUT requests.
func (u userController) UpdateUser(c *gin.Context) {
	userService := u.userService

	var p types.UserUpdateInput
	if err := c.BindJSON(&p); err != nil {
		return
	}

	var oldUser types.UserLoginInput
	oldUser.UserName, _ = c.Params.Get("userName")
	oldUser.Password = p.OldPassword

	var newUser types.UserLoginInput
	newUser.UserName = oldUser.UserName
	newUser.Password = p.NewPassword

	user, err := userService.UpdateUser(&oldUser, &newUser)
	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusOK, user)

	case errortypes.IncorrectUsernameOrPasswordError:
		_ = c.AbortWithError(http.StatusUnauthorized, err)

	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedUserError{User: types.User{UserName: oldUser.UserName}})
	}
}
