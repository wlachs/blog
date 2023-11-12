package errortypes

import (
	"fmt"
	"github.com/wlchs/blog/internal/types"
)

type UnexpectedPostError struct {
	Post types.Post
}

func (e UnexpectedPostError) Error() string {
	if e.Post.URLHandle != "" {
		return fmt.Sprintf("unexpected error encountered with post \"%s\"", e.Post.URLHandle)
	}
	return "unexpected post error encountered"
}

type MissingUrlHandleError struct{}

func (e MissingUrlHandleError) Error() string {
	return "no URL handle provided"
}

type PostNotFoundError struct {
	Post types.Post
}

func (e PostNotFoundError) Error() string {
	return fmt.Sprintf("post with URL handle \"%s\" not found", e.Post.URLHandle)
}
