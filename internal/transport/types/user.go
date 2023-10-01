package types

type UserLoginInput struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserPasswordChangeInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type User struct {
	UserName     string   `json:"userName"`
	PasswordHash string   `json:"-"`
	Posts        []string `json:"posts"`
}
