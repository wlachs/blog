package errors

import "fmt"

type IncorrectSecretError struct{}

func (i IncorrectSecretError) Error() string {
	return fmt.Sprintln("incorrect registration secret")
}
