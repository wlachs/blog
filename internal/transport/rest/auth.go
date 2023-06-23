package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
	"github.com/wlchs/blog/internal/transport/types"
)

func loginMiddleware(c *gin.Context) {
	var u types.UserLoginInput

	if err := c.BindJSON(&u); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}

	token, err := services.AuthenticateUser(u)

	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.Header("X-Auth-Token", token)
	c.Status(http.StatusOK)
}

func jwtAuthMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("X-Auth-Token")

	if token == "" {
		c.AbortWithStatus(401)
	} else if services.VerifyAuthenticationToken(token) {
		c.Next()
	} else {
		c.AbortWithStatus(401)
	}
}
