package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString takes a string as input and calculates its hash.
func HashString(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), 10)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareStringWithHash takes a plaintext string and a hash as input and compares the hash of the plaintext to the provided hash.
func CompareStringWithHash(s string, h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(s)) == nil
}
