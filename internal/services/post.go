package services

import (
	"github.com/wlchs/blog/internal/container"
	"github.com/wlchs/blog/internal/repository"
	"github.com/wlchs/blog/internal/types"
)

// PostService interface. Defines post-related business logic.
type PostService interface {
	AddPost(newPost *types.Post) (types.Post, error)
	GetPost(id string) (types.Post, error)
	GetPosts() ([]types.Post, error)
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
func (p postService) AddPost(newPost *types.Post) (types.Post, error) {
	log := p.cont.GetLogger()
	postRepository := p.cont.GetPostRepository()
	userRepository := p.cont.GetUserRepository()

	// Get post author
	author, err := userRepository.GetUser(newPost.Author)
	if err != nil {
		log.Errorf("failed to get author for post %v with username %s", newPost, newPost.Author)
		return types.Post{}, err
	}

	log.Infof("adding new post %v with author %s", newPost, newPost.Author)

	post, err := postRepository.AddPost(newPost, author.ID)
	return mapPost(post), err
}

// GetPost retrieves the post with the given URL handle.
func (p postService) GetPost(urlHandle string) (types.Post, error) {
	postRepository := p.cont.GetPostRepository()
	post, err := postRepository.GetPost(urlHandle)
	return mapPost(post), err
}

// GetPosts retrieves every post of the blog.
func (p postService) GetPosts() ([]types.Post, error) {
	postRepository := p.cont.GetPostRepository()
	posts, err := postRepository.GetPosts()
	return mapPosts(posts), err
}

// mapPost maps a Post model to a post data object
func mapPost(p *repository.Post) types.Post {
	if p == nil {
		return types.Post{}
	}
	return types.Post{
		URLHandle:    p.URLHandle,
		Title:        p.Title,
		Author:       p.Author.UserName,
		Summary:      p.Summary,
		Body:         p.Body,
		CreationTime: p.CreatedAt,
	}
}

// mapPostMetadata maps a Post model to a post metadata object
func mapPostMetadata(p *repository.Post) types.Post {
	return types.Post{
		URLHandle:    p.URLHandle,
		Title:        p.Title,
		Author:       p.Author.UserName,
		Summary:      p.Summary,
		CreationTime: p.CreatedAt,
	}
}

// mapPosts maps a slice of Post models to a slice of post data objects
func mapPosts(p []repository.Post) []types.Post {
	if p == nil {
		return []types.Post{}
	}
	posts := make([]types.Post, 0, len(p))

	for _, post := range p {
		posts = append(posts, mapPostMetadata(&post))
	}

	return posts
}

// mapPostHandles maps a slice of Post models to a slice of strings containing the posts' URL handles
func mapPostHandles(p []repository.Post) []string {
	handles := make([]string, 0, len(p))

	for _, post := range p {
		handles = append(handles, post.URLHandle)
	}

	return handles
}
