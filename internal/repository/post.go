package repository

//go:generate mockgen-v0.4.0 -source=post.go -destination=../mocks/mock_post_repository.go -package=mocks

import (
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/wlachs/blog/internal/errortypes"
)

// Post DB schema
type Post struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	URLHandle string `gorm:"unique;not null"`
	AuthorID  uint   `gorm:"not null"`
	Author    User
	Title     string
	Summary   string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PostRepository interface defining post-related database operations.
type PostRepository interface {
	AddPost(post Post, authorID uint) (Post, error)
	GetPost(urlHandle string) (Post, error)
	GetPosts() ([]Post, error)
}

// postRepository is the concrete implementation of the PostRepository interface.
type postRepository struct {
	logger     *zap.SugaredLogger
	repository Repository
}

// CreatePostRepository instantiates the postRepository
func CreatePostRepository(logger *zap.SugaredLogger, repository Repository) PostRepository {
	initPostModel(logger, repository)

	return &postRepository{
		logger:     logger,
		repository: repository,
	}
}

// initPostModel initializes the Post schema in the database
func initPostModel(logger *zap.SugaredLogger, repository Repository) {
	if err := repository.AutoMigrate(&Post{}); err != nil {
		logger.Errorf("failed to initialize post model: %v", err)
	}
}

// AddPost adds a new post with the provided fields to the database.
// The second parameter holds information about the author.
func (p postRepository) AddPost(post Post, authorID uint) (Post, error) {
	log := p.logger
	repo := p.repository

	newPost := Post{
		URLHandle: post.URLHandle,
		Title:     post.Title,
		Summary:   post.Summary,
		Body:      post.Body,
		AuthorID:  authorID,
	}

	if result := repo.Create(&newPost); result.Error == nil {
		log.Debugf("created post: %v", newPost)
		return newPost, nil
	} else if strings.Contains(result.Error.Error(), "1062") {
		log.Debugf("failed to create post, duplicate key: %s, error: %v", post.URLHandle, result.Error)
		return Post{}, errortypes.DuplicateElementError{Key: post.URLHandle}
	} else {
		log.Debugf("failed to create post: %v, error: %s", post, result.Error)
		return Post{}, result.Error
	}
}

// GetPost retrieves the post with the given URL-handle from the database.
func (p postRepository) GetPost(urlHandle string) (Post, error) {
	log := p.logger
	repo := p.repository

	post := Post{
		URLHandle: urlHandle,
	}

	result := repo.Preload("Author").Where(&post).Take(&post)

	if result.Error != nil {
		log.Debugf("failed to retrieve post with handle: %s, error: %v", urlHandle, result.Error)
		if result.Error.Error() == "record not found" {
			return Post{}, errortypes.PostNotFoundError{URLHandle: urlHandle}
		}
		return Post{}, result.Error
	}

	log.Debugf("retrieved post: %v", post)
	return post, nil
}

// GetPosts retrieves every post from the database.
func (p postRepository) GetPosts() ([]Post, error) {
	log := p.logger
	repo := p.repository

	var posts []Post
	if result := repo.Preload("Author").Order("created_at DESC").Find(&posts); result.Error != nil {
		log.Debugf("error fetching posts: %v", result.Error)
		return []Post{}, result.Error
	}

	log.Debugf("fetched posts: %v", posts)
	return posts, nil
}
