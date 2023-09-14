package types

type UserLoginInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type User struct {
	UserName     string   `json:"userName"`
	PasswordHash string   `json:"-"`
	Posts        []string `json:"posts"`
}
