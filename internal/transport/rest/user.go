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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newUser, err := services.RegisterUser(u)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, newUser)

	case errors.IncorrectSecretError:
		c.AbortWithError(http.StatusUnauthorized, err)

	default:
		c.AbortWithError(http.StatusBadRequest, err)
	}
}
