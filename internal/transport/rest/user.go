package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/transport/types"
)

func registerHandler(c *gin.Context) {
	var u types.UserInput

	if err := c.BindJSON(&u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newUser, err := services.RegisterUser(u)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}
