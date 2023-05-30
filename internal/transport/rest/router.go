package rest

import (
	"os"

	"github.com/gin-gonic/gin"
)

func InitRoutes() error {
	router := gin.Default()

	// Posts
	router.GET("/posts", getPostsMiddleware)
	router.GET("/posts/:id", getPostMiddleware)
	router.POST("/posts", jwtAuthMiddleware, addPostMiddleware)

	// Users
	router.GET("/users", getUsersMiddleware)
	router.GET("/users/:userName", getUserMiddleware)
	router.POST("/login", loginMiddleware)
	router.POST("/register", registerHandler)

	PORT := os.Getenv("PORT")
	return router.Run(":" + PORT)
}
