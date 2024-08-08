package errortypes

import (
	"fmt"
)

type UnexpectedPostError struct {
	URLHandle string
}

func (e UnexpectedPostError) Error() string {
	if e.URLHandle != "" {
		return fmt.Sprintf("unexpected error encountered with post \"%s\"", e.URLHandle)
	}
	return "unexpected post error encountered"
}

type MissingUrlHandleError struct{}

func (e MissingUrlHandleError) Error() string {
	return "no URL handle provided"
}

type PostNotFoundError struct {
	URLHandle string
}

func (e PostNotFoundError) Error() string {
	return fmt.Sprintf("post with URL handle \"%s\" not found", e.URLHandle)
}

type InvalidPostPageError struct {
	Page int
}

func (e InvalidPostPageError) Error() string {
	return fmt.Sprintf("post page with number %d not valid", e.Page)
}
