package services

//go:generate mockgen-v0.4.0 -source=post.go -destination=../mocks/mock_post_service.go -package=mocks

import (
	"github.com/wlachs/blog/internal/container"
	"github.com/wlachs/blog/internal/repository"
)

// PostService interface. Defines post-related business logic.
type PostService interface {
	AddPost(newPost repository.Post, authorName string) (repository.Post, error)
	UpdatePost(updatedPost repository.Post) (repository.Post, error)
	DeletePost(id string) error
	GetPost(id string) (repository.Post, error)
	GetPosts() ([]repository.Post, error)
	GetPostsPage(page int) ([]repository.Post, error)
}

// postService is the concrete implementation of the PostService interface.
type postService struct {
	cont container.Container
}

// CreatePostService instantiates the postService using the application container.
func CreatePostService(cont container.Container) PostService {
	return &postService{cont}
}

// AddPost adds a new post to the blog.
func (p postService) AddPost(newPost repository.Post, authorName string) (repository.Post, error) {
	log := p.cont.GetLogger()
	postRepository := p.cont.GetPostRepository()
	userRepository := p.cont.GetUserRepository()

	// Get post author
	author, err := userRepository.GetUser(authorName)
	if err != nil {
		log.Errorf("failed to get author for post %v with username %s", newPost, authorName)
		return repository.Post{}, err
	}

	newPost.AuthorID = author.ID

	log.Infof("adding new post %v with author %s", newPost, authorName)
	return postRepository.AddPost(newPost)
}

// UpdatePost updates an existing post in the blog.
func (p postService) UpdatePost(updatedPost repository.Post) (repository.Post, error) {
	log := p.cont.GetLogger()
	postRepository := p.cont.GetPostRepository()

	log.Infof("updating post %v", updatedPost)
	return postRepository.UpdatePost(updatedPost)
}

// DeletePost deletes a post from the blog.
func (p postService) DeletePost(urlHandle string) error {
	log := p.cont.GetLogger()
	postRepository := p.cont.GetPostRepository()

	log.Infof("deleting post %s", urlHandle)
	return postRepository.DeletePost(urlHandle)
}

// GetPost retrieves the post with the given URL handle.
func (p postService) GetPost(urlHandle string) (repository.Post, error) {
	postRepository := p.cont.GetPostRepository()
	return postRepository.GetPost(urlHandle)
}

// GetPosts retrieves the first page of posts of the blog.
func (p postService) GetPosts() ([]repository.Post, error) {
	return p.GetPostsPage(1)
}

// GetPostsPage retrieves one page of posts of the blog.
func (p postService) GetPostsPage(page int) ([]repository.Post, error) {
	postRepository := p.cont.GetPostRepository()
	return postRepository.GetPosts()
}
