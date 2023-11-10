package controller

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
)

// AuthController interface defining authentication-related methods to handler HTTP requests.
type AuthController interface {
	Login(c *gin.Context)
	Protect(c *gin.Context)
}

// authController is a concrete implementation of the AuthController interface.
type authController struct {
	cont container.Container
}

// CreateAuthController instantiates the AuthController using the application container.
func CreateAuthController(cont container.Container) AuthController {
	return &authController{cont: cont}
}

// Login middleware. Top level handler of /login POST requests.
func (auth authController) Login(c *gin.Context) {
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

// Protect middleware. Can be used before any middleware to make sure only authenticated users are able to use an endpoint.
func (auth authController) Protect(c *gin.Context) {
	token := c.Request.Header.Get("X-Auth-Token")

	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else if u, err := services.ParseJWT(token); err == nil {
		c.Set("user", u)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
