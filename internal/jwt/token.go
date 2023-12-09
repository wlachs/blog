package jwt

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// signingKey is the JWT secret key stored as an environment variable
var signingKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

// TokenUtils interface. JWT-related utility methods.
type TokenUtils interface {
	ParseJWT(t string) (string, error)
	GenerateJWT(userName string) (string, error)
}

// tokenUtils struct. Placeholder receiver struct for JWT utils.
type tokenUtils struct {
	logger *zap.SugaredLogger
}

// CreateTokenUtils instantiates the tokenUtils implementation.
func CreateTokenUtils(logger *zap.SugaredLogger) TokenUtils {
	return &tokenUtils{
		logger: logger,
	}
}

// ParseJWT parses a token and extracts the user field if valid.
func (j tokenUtils) ParseJWT(t string) (string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user"].(string), nil
	} else {
		return "", fmt.Errorf("failed to get jwt claims")
	}
}

// GenerateJWT creates a JWT containing the following fields:
// - username
// - authorized flag
// - expiration date
func (j tokenUtils) GenerateJWT(userName string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = userName

	return token.SignedString(signingKey)
}
