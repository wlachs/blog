package types

type UserLoginInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserRegisterInput struct {
	UserLoginInput
	RegistrationSecret string `json:"registrationSecret"`
}

type User struct {
	UserName     string `json:"userName"`
	PasswordHash string `json:"-"`
}
