package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wlachs/blog/api/types"
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/repository"
	"github.com/wlachs/blog/internal/services"
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

	var body types.NewPost
	if err := c.BindJSON(&body); err != nil {
		return
	}

	// Set author and post ID from context
	author := c.GetString("UserID")
	postID, _ := c.Params.Get("PostID")

	// Create new raw post item
	newPost := repository.Post{
		URLHandle: postID,
	}
	if body.Title != nil {
		newPost.Title = *body.Title
	}
	if body.Summary != nil {
		newPost.Summary = *body.Summary
	}
	if body.Body != nil {
		newPost.Body = *body.Body
	}

	post, err := postService.AddPost(newPost, author)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusCreated, populatePost(post))
	case errortypes.DuplicateElementError:
		_ = c.AbortWithError(http.StatusConflict, err)
	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedPostError{URLHandle: postID})
	}
}

// GetPost middleware. Top level handler of /posts/:id GET requests.
func (controller postController) GetPost(c *gin.Context) {
	postService := controller.postService

	id, found := c.Params.Get("PostID")
	if !found {
		_ = c.AbortWithError(http.StatusBadRequest, errortypes.MissingUrlHandleError{})
		return
	}

	post, err := postService.GetPost(id)

	switch err.(type) {
	case nil:
		c.IndentedJSON(http.StatusOK, populatePost(post))
	case errortypes.PostNotFoundError:
		_ = c.AbortWithError(http.StatusNotFound, err)
	default:
		_ = c.AbortWithError(http.StatusInternalServerError, errortypes.UnexpectedPostError{URLHandle: id})
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

	c.IndentedJSON(http.StatusOK, populatePostMetadataSlice(posts))
}

// populatePost maps a repository.Post model to types.Post
func populatePost(post repository.Post) types.Post {
	p := types.Post{
		Author:       post.Author.UserName,
		CreationTime: post.CreatedAt,
		Id:           post.URLHandle,
		Title:        post.Title,
	}

	if len(post.Summary) > 0 {
		p.Summary = &post.Summary
	}
	if len(post.Body) > 0 {
		p.Body = &post.Body
	}

	return p
}

// populatePostMetadata maps a repository.Post model to types.PostMetadata
func populatePostMetadata(post repository.Post) types.PostMetadata {
	p := types.PostMetadata{
		Author:       post.Author.UserName,
		CreationTime: post.CreatedAt,
		Id:           post.URLHandle,
		Title:        post.Title,
	}

	if len(post.Summary) > 0 {
		p.Summary = &post.Summary
	}

	return p
}

// populatePostMetadataSlice maps a slice of repository.Post models to a types.PostMetadata slice
func populatePostMetadataSlice(posts []repository.Post) []types.PostMetadata {
	p := make([]types.PostMetadata, 0, len(posts))

	for _, post := range posts {
		p = append(p, populatePostMetadata(post))
	}

	return p
}
