package services

import (
	"github.com/wlchs/blog/internal/models"
	"github.com/wlchs/blog/internal/transport/types"
)

func mapPost(p models.Post) types.Post {
	return types.Post{
		URLHandle:    p.URLHandle,
		Title:        p.Title,
		Author:       p.Author.UserName,
		Summary:      p.Summary,
		Body:         p.Body,
		CreationTime: p.CreatedAt,
	}
}

func mapPosts(p []models.Post) []types.Post {
	var posts []types.Post

	for _, post := range p {
		posts = append(posts, mapPost(post))
	}

	return posts
}

func mapPostHandles(p []models.Post) []string {
	var handles []string

	for _, post := range p {
		handles = append(handles, post.URLHandle)
	}

	return handles
}

func GetPosts() ([]types.Post, error) {
	p, err := models.GetPosts()
	return mapPosts(p), err
}

func GetPost(id string) (types.Post, error) {
	p, err := models.GetPost(id)
	return mapPost(p), err
}

func AddPost(post types.Post) (types.Post, error) {
	p, err := models.AddPost(post)
	return mapPost(p), err
}
