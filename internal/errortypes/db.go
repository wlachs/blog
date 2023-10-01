package errortypes

import "fmt"

type DuplicateElementError struct {
	Key string
}

func (d DuplicateElementError) Error() string {
	return fmt.Sprintf("object with key \"%s\" already exists", d.Key)
}
