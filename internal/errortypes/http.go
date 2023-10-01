package errortypes

import "fmt"

type ErrorWithStatus struct {
	Status int
	Err    error
}

func (e ErrorWithStatus) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Err.Error())
}
