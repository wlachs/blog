package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/services"
	"os"
)

// CreateRoutes initializes and serves the REST API
func CreateRoutes(cont container.Container) {
	log := cont.GetLogger()
	router := gin.Default()

	// Services
	postService := services.CreatePostService(cont)
	userService := services.CreateUserService(cont)

	// Controllers
	authCtrl := CreateAuthController(cont, userService)
	postCtrl := CreatePostController(cont, postService)
	userCtrl := CreateUserController(cont, userService)

	// Posts
	router.GET("/posts", postCtrl.GetPosts)
	router.GET("/posts/:id", postCtrl.GetPost)
	router.POST("/posts", authCtrl.Protect, postCtrl.AddPost)

	// Users
	router.GET("/users", userCtrl.GetUsers)
	router.GET("/users/:userName", userCtrl.GetUser)
	router.PUT("/users/:userName", userCtrl.UpdateUser)
	router.POST("/login", authCtrl.Login)

	port := os.Getenv("PORT")
	err := router.Run(":" + port)

	if err != nil {
		log.Errorf("error encountered in router: %v", err)
		os.Exit(1)
	}
}
