package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/services"
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
	router.GET("/api/v0/posts", postCtrl.GetPosts)
	router.GET("/api/v0/posts/:PostID", postCtrl.GetPost)
	router.POST("/api/v0/posts/:PostID", authCtrl.Protect, postCtrl.AddPost)
	router.PUT("/api/v0/posts/:PostID", authCtrl.Protect, postCtrl.UpdatePost)
	router.DELETE("/api/v0/posts/:PostID", authCtrl.Protect, postCtrl.DeletePost)

	// Users
	router.GET("/api/v0/users", userCtrl.GetUsers)
	router.GET("/api/v0/users/:UserID", userCtrl.GetUser)
	router.POST("/api/v0/users/:UserID", authCtrl.Protect, userCtrl.AddUser)
	router.PUT("/api/v0/users/:UserID", authCtrl.Protect, userCtrl.UpdateUser)
	router.DELETE("/api/v0/users/:UserID", authCtrl.Protect, userCtrl.DeleteUser)
	router.POST("/api/v0/login", authCtrl.Login)

	port := os.Getenv("PORT")
	err := router.Run(":" + port)

	if err != nil {
		log.Errorf("error encountered in router: %v", err)
	}
}
