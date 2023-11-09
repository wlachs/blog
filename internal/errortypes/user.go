package errortypes

type IncorrectUsernameOrPasswordError struct{}

func (i IncorrectUsernameOrPasswordError) Error() string {
	return "incorrect username or password"
}
