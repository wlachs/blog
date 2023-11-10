package errortypes

type IncorrectUsernameOrPasswordError struct{}
type PasswordHashingError struct{}

func (i IncorrectUsernameOrPasswordError) Error() string {
	return "incorrect username or password"
}

func (i PasswordHashingError) Error() string {
	return "password hashing failed"
}
