package errortypes

import (
	"fmt"
)

type IncorrectUsernameOrPasswordError struct{}

func (i IncorrectUsernameOrPasswordError) Error() string {
	return "incorrect username or password"
}

type MissingDefaultUsernameOrPasswordError struct{}

func (i MissingDefaultUsernameOrPasswordError) Error() string {
	return "missing default username or password"
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

type MissingPasswordError struct {
}

func (i MissingPasswordError) Error() string {
	return "no password provided"
}

type UserNotFoundError struct {
	UserName string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user \"%s\" not found", e.UserName)
}

type UnexpectedUserError struct {
	UserName string
}

func (e UnexpectedUserError) Error() string {
	if e.UserName != "" {
		return fmt.Sprintf("unexpected error encountered with user \"%s\"", e.UserName)
	}
	return "unexpected user error encountered"
}
