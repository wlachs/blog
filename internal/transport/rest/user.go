package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/errors"
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/transport/types"
)

func registerHandler(c *gin.Context) {
	var u types.UserRegisterInput

	if err := c.BindJSON(&u); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newUser, err := services.RegisterUser(u)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, newUser)

	case errors.IncorrectSecretError:
		_ = c.AbortWithError(http.StatusUnauthorized, err)

	default:
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
}

func getUsersMiddleware(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func getUserMiddleware(c *gin.Context) {
	userName, found := c.Params.Get("userName")

	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No userName provided!")
		return
	}

	user, err := services.GetUser(userName)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
