package controller

import (
	"github.com/wlchs/blog/internal/container"
	"os"

	"github.com/gin-gonic/gin"
)

// CreateRoutes initializes and serves the REST API
func CreateRoutes(cont container.Container) {
	log := cont.GetLogger()
	router := gin.Default()

	// Controllers
	auth := CreateAuthController(cont)
	post := CreatePostController(cont)
	user := CreateUserController(cont)

	// Posts
	router.GET("/posts", post.GetPosts)
	router.GET("/posts/:id", post.GetPost)
	router.POST("/posts", auth.Protect, post.AddPost)

	// Users
	router.GET("/users", user.GetUsers)
	router.GET("/users/:userName", user.GetUser)
	router.PUT("/users/:userName", user.UpdateUser)
	router.POST("/login", auth.Login)

	port := os.Getenv("PORT")
	err := router.Run(":" + port)

	if err != nil {
		log.Error("error encountered in router: %s", err)
		os.Exit(1)
	}
}
