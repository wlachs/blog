package errortypes

import (
	"fmt"
	"github.com/wlchs/blog/internal/types"
)

type IncorrectUsernameOrPasswordError struct{}

func (i IncorrectUsernameOrPasswordError) Error() string {
	return "incorrect username or password"
}

type PasswordHashingError struct{}

func (i PasswordHashingError) Error() string {
	return "password hashing failed"
}

type MissingUsernameError struct {
}

func (i MissingUsernameError) Error() string {
	return "no username provided"
}

type UserNotFoundError struct {
	User types.User
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user \"%s\" not found", e.User.UserName)
}

type UnexpectedUserError struct {
	User types.User
}

func (e UnexpectedUserError) Error() string {
	if e.User.UserName != "" {
		return fmt.Sprintf("unexpected error encountered with user \"%s\"", e.User.UserName)
	}
	return "unexpected user error encountered"
}
