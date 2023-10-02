package rest

import (
	"errors"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/transport/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
)

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

func updateUserMiddleware(c *gin.Context) {
	var p types.UserUpdateInput
	if err := c.BindJSON(&p); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var oldUser types.UserLoginInput
	oldUser.UserName, _ = c.Params.Get("userName")
	oldUser.Password = p.OldPassword

	hash, err := services.HashString(p.NewPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	var newUser types.User
	newUser.UserName = oldUser.UserName
	newUser.PasswordHash = hash

	user, err := services.UpdateUser(oldUser, newUser)
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

	c.IndentedJSON(http.StatusOK, user)
}
