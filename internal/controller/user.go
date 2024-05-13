package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wlachs/blog/api/types"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/repository"
	"github.com/wlachs/blog/internal/services"
	"net/http"
)

// UserController interface defining user-related middleware methods to handler HTTP requests.
type UserController interface {
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
	AddUser(c *gin.Context)
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

// GetUser middleware. Top level handler of /user/:UserID GET requests.
func (u userController) GetUser(c *gin.Context) {
	userService := u.userService
	userName, found := c.Params.Get("UserID")

	if !found {
		_ = c.AbortWithError(http.StatusBadRequest, errortypes.MissingUsernameError{})
		return
	}

	user, err := userService.GetUser(userName)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusOK, populateUser(user))
	case errortypes.UserNotFoundError:
		_ = c.AbortWithError(http.StatusNotFound, err)
	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedUserError{UserName: userName})
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

	c.IndentedJSON(http.StatusOK, populateUsers(users))
}

// AddUser middleware. Top level handler of /users/:UserID POST requests.
// Registers a new user.
func (u userController) AddUser(c *gin.Context) {
	userService := u.userService

	var p types.AddUserJSONBody
	if err := c.BindJSON(&p); err != nil {
		return
	}

	userID, _ := c.Params.Get("UserID")
	user, err := userService.RegisterUser(userID, p.Password)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, populateUser(user))
	case errortypes.MissingPasswordError:
		_ = c.AbortWithError(http.StatusBadRequest, err)
	case errortypes.DuplicateElementError:
		_ = c.AbortWithError(http.StatusConflict, err)
	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedUserError{UserName: userID})
	}
}

// UpdateUser middleware. Top level handler of /users/:UserID PUT requests.
func (u userController) UpdateUser(c *gin.Context) {
	userService := u.userService

	var p types.UpdateUserJSONBody
	if err := c.BindJSON(&p); err != nil {
		return
	}

	userID, _ := c.Params.Get("UserID")
	user, err := userService.UpdateUser(userID, p.OldPassword, p.NewPassword)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusOK, populateUser(user))
	case errortypes.IncorrectUsernameOrPasswordError:
		_ = c.AbortWithError(http.StatusUnauthorized, err)
	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedUserError{UserName: userID})
	}
}

// populateUser maps a repository.User model to types.User
func populateUser(user repository.User) types.User {
	u := types.User{
		UserID: user.UserName,
	}

	if len(user.Posts) > 0 {
		posts := populatePostMetadataSlice(user.Posts)
		u.Posts = &posts
	}

	return u
}

// populateUsers maps a slice of repository.User models to a types.User slice
func populateUsers(users []repository.User) []types.User {
	u := make([]types.User, 0, len(users))

	for _, user := range users {
		u = append(u, populateUser(user))
	}

	return u
}
