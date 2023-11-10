package errortypes

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
