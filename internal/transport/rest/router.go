package rest

import (
	"github.com/wlchs/blog/internal/container"
	"os"

	"github.com/gin-gonic/gin"
)

// CreateRoutes initializes and serves the REST API
func CreateRoutes(cont container.Container) {
	log := cont.GetLogger()
	router := gin.Default()

	// Posts
	router.GET("/posts", getPostsMiddleware)
	router.GET("/posts/:id", getPostMiddleware)
	router.POST("/posts", jwtAuthMiddleware, addPostMiddleware)

	// Users
	router.GET("/users", getUsersMiddleware)
	router.GET("/users/:userName", getUserMiddleware)
	router.PUT("/users/:userName", updateUserMiddleware)
	router.POST("/login", loginMiddleware)

	port := os.Getenv("PORT")
	err := router.Run(":" + port)

	if err != nil {
		log.Error("error encountered in router: %s", err)
		os.Exit(1)
	}
}
