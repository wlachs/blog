package controller

import (
	"errors"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
)

// UserController interface defining user-related middleware methods to handler HTTP requests.
type UserController interface {
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
}

// userController is a concrete implementation of the UserController interface.
type userController struct {
	cont container.Container
}

// CreateUserController instantiates a user controller user the application container.
func CreateUserController(cont container.Container) UserController {
	return &userController{cont: cont}
}

// GetUser middleware. Top level handler of /user/:userName GET requests.
func (user userController) GetUser(c *gin.Context) {
	userName, found := c.Params.Get("userName")

	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No userName provided!")
		return
	}

	u, err := services.GetUser(userName)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, u)
}

// GetUsers middleware. Top level handler of /users GET requests.
func (user userController) GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// UpdateUser middleware. Top level handler of /users/:userName PUT requests.
func (user userController) UpdateUser(c *gin.Context) {
	var p types.UserUpdateInput
	if err := c.BindJSON(&p); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var oldUser types.UserLoginInput
	oldUser.UserName, _ = c.Params.Get("userName")
	oldUser.Password = p.OldPassword

	var newUser types.UserLoginInput
	newUser.UserName = oldUser.UserName
	newUser.Password = p.NewPassword

	u, err := services.UpdateUser(oldUser, newUser)
	if err != nil {
		var incorrectUsernameOrPasswordError errortypes.IncorrectUsernameOrPasswordError
		if errors.As(err, &incorrectUsernameOrPasswordError) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, incorrectUsernameOrPasswordError.Error())
			return
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, u)
}
