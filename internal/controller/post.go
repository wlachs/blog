package controller

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wlchs/blog/internal/services"
)

// PostController interface defining post-related middleware methods to handle HTTP requests
type PostController interface {
	AddPost(c *gin.Context)
	GetPost(c *gin.Context)
	GetPosts(c *gin.Context)
}

// postController is a concrete implementation of the PostController interface
type postController struct {
	cont        container.Container
	postService services.PostService
}

// CreatePostController instantiates a post controller using the application container.
func CreatePostController(cont container.Container) PostController {
	return &postController{cont, services.CreatePostService(cont)}
}

// AddPost middleware. Top level handler of /posts POST requests.
func (controller postController) AddPost(c *gin.Context) {
	postService := controller.postService

	var post types.Post
	if err := c.BindJSON(&post); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Set author from context
	post.Author = c.GetString("user")
	post, err := postService.AddPost(&post)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, post)

	case errortypes.DuplicateElementError:
		_ = c.AbortWithError(http.StatusConflict, err)

	default:
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
}

// GetPost middleware. Top level handler of /posts/:id GET requests.
func (controller postController) GetPost(c *gin.Context) {
	postService := controller.postService

	id, found := c.Params.Get("id")
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, "No id provided!")
		return
	}

	post, err := postService.GetPost(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

// GetPosts middleware. Top level handler of /posts GET requests.
func (controller postController) GetPosts(c *gin.Context) {
	postService := controller.postService

	posts, err := postService.GetPosts()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}
