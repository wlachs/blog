package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/services"
	"github.com/wlachs/blog/internal/types"
	"net/http"
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
func CreatePostController(cont container.Container, postService services.PostService) PostController {
	return &postController{cont, postService}
}

// AddPost middleware. Top level handler of /posts POST requests.
func (controller postController) AddPost(c *gin.Context) {
	postService := controller.postService

	var body types.Post
	if err := c.BindJSON(&body); err != nil {
		return
	}

	// Set author from context
	body.Author = c.GetString("user")
	post, err := postService.AddPost(&body)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, post)

	case errortypes.DuplicateElementError:
		_ = c.AbortWithError(http.StatusConflict, err)

	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedPostError{Post: body})
	}
}

// GetPost middleware. Top level handler of /posts/:id GET requests.
func (controller postController) GetPost(c *gin.Context) {
	postService := controller.postService

	id, found := c.Params.Get("id")
	if !found {
		_ = c.AbortWithError(http.StatusBadRequest, errortypes.MissingUrlHandleError{})
		return
	}

	post, err := postService.GetPost(id)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusOK, post)

	case errortypes.PostNotFoundError:
		_ = c.AbortWithError(http.StatusNotFound, err)

	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedPostError{Post: types.Post{URLHandle: id}})
	}
}

// GetPosts middleware. Top level handler of /posts GET requests.
func (controller postController) GetPosts(c *gin.Context) {
	postService := controller.postService

	posts, err := postService.GetPosts()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedPostError{})
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}
