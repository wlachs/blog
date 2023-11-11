package errortypes

type MissingAuthTokenError struct{}

func (m MissingAuthTokenError) Error() string {
	return "auth token is missing"
}

type InvalidAuthTokenError struct{}

func (i InvalidAuthTokenError) Error() string {
	return "auth token expired or invalid"
}
