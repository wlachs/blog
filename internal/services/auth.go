package services

import (
	"fmt"

	"github.com/wlchs/blog/internal/transport/types"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(u types.UserInput) (string, error) {
	if !CheckUserPassword(u) {
		return "", fmt.Errorf("incorrect username or password")
	}

	return GenerateJWT(u.UserName)
}

func VerifyAuthenticationToken(t string) bool {
	_, err := ParseJWT(t)
	return err == nil
}

func CheckUserPassword(u types.UserInput) bool {
	user, _ := GetUser(u.UserName)
	return CompareStringWithHash(u.Password, user.PasswordHash)
}

func HashString(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), 10)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareStringWithHash(s string, h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(s)) == nil
}
