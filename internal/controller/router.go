package controller

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/services"
	"os"

	"github.com/gin-gonic/gin"
)

// CreateRoutes initializes and serves the REST API
func CreateRoutes(cont container.Container) {
	log := cont.GetLogger()
	router := gin.Default()

	// Services
	postService := services.CreatePostService(cont)
	userService := services.CreateUserService(cont)

	// Controllers
	authController := CreateAuthController(cont, userService)
	postController := CreatePostController(cont, postService)
	userController := CreateUserController(cont, userService)

	// Posts
	router.GET("/posts", postController.GetPosts)
	router.GET("/posts/:id", postController.GetPost)
	router.POST("/posts", authController.Protect, postController.AddPost)

	// Users
	router.GET("/users", userController.GetUsers)
	router.GET("/users/:userName", userController.GetUser)
	router.PUT("/users/:userName", userController.UpdateUser)
	router.POST("/login", authController.Login)

	port := os.Getenv("PORT")
	err := router.Run(":" + port)

	if err != nil {
		log.Errorf("error encountered in router: %v", err)
		os.Exit(1)
	}
}
