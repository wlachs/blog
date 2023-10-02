package errortypes

import "fmt"

type IncorrectUsernameOrPasswordError struct{}

func (i IncorrectUsernameOrPasswordError) Error() string {
	return fmt.Sprintf("incorrect username or password")
}
