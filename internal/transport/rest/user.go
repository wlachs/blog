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

func passwordChangeMiddleware(c *gin.Context) {
	var p types.UserPasswordChangeInput
	if err := c.BindJSON(&p); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var oldUser types.UserLoginInput
	oldUser.UserName = c.GetString("user")
	oldUser.Password = p.OldPassword

	var newUser types.UserLoginInput
	newUser.UserName = oldUser.UserName
	newUser.Password = p.NewPassword

	user, err := services.ChangeUserPassword(oldUser, newUser)
	if err != nil {
		var errorWithStatus errortypes.ErrorWithStatus
		if errors.As(err, &errorWithStatus) {
			c.AbortWithStatusJSON(errorWithStatus.Status, errorWithStatus.Error())
			return
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, user)
}
