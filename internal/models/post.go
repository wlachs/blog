package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/wlchs/blog/internal/database"
	"github.com/wlchs/blog/internal/errortypes"
	"github.com/wlchs/blog/internal/transport/types"
)

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

func GetPosts() ([]Post, error) {
	var p []Post
	if result := database.Agent.Preload("Author").Order("created_at DESC").Find(&p); result.Error != nil {
		return []Post{}, result.Error
	}

	return p, nil
}

func GetPost(h string) (Post, error) {
	p := Post{
		URLHandle: h,
	}

	result := database.Agent.Preload("Author").Where(&p).Take(&p)

	if result.Error != nil {
		return Post{}, result.Error
	}

	if result.RowsAffected == 0 {
		return Post{}, fmt.Errorf("post with handle: %s not found", h)
	}

	return p, nil
}

func AddPost(p types.Post) (Post, error) {
	u, err := GetUser(p.Author)
	if err != nil {
		return Post{}, err
	}

	newPost := Post{
		URLHandle: p.URLHandle,
		Author:    u,
		Title:     p.Title,
		Summary:   p.Summary,
		Body:      p.Body,
	}
	if result := database.Agent.Create(&newPost); result.Error == nil {
		return newPost, nil
	} else if strings.Contains(result.Error.Error(), "1062") {
		return Post{}, errortypes.DuplicateElementError{Key: p.URLHandle}
	}

	return Post{}, fmt.Errorf("error creating post with URL handle \"%s\"", newPost.URLHandle)
}
