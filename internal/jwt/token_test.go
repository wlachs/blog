package jwt_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/jwt"
	"github.com/wlchs/blog/internal/logger"
	"testing"
)

// tokenUtilsTestContext contains objects relevant for testing the TokenUtils.
type tokenUtilsTestContext struct {
	sut jwt.TokenUtils
}

// createTokenUtilsContext creates the context for testing the TokenUtils and reduces code duplication.
func createTokenUtilsContext(t *testing.T) *tokenUtilsTestContext {
	t.Helper()

	sut := jwt.CreateTokenUtils(logger.CreateLogger())

	return &tokenUtilsTestContext{sut}
}

// TestTokenUtils_GenerateJWT tests generating a new token
func TestTokenUtils_GenerateJWT(t *testing.T) {
	t.Parallel()
	c := createTokenUtilsContext(t)

	userName := "TestAuthor"

	token, err := c.sut.GenerateJWT(userName)
	assert.Greater(t, len(token), 0, "token shouldn't be empty")
	assert.Nil(t, err, "expected to complete without error")
}

// TestTokenUtils_ParseJWT tests parsing a valid JWT
func TestTokenUtils_ParseJWT(t *testing.T) {
	t.Parallel()
	c := createTokenUtilsContext(t)

	expectedUserName := "TestAuthor"

	token, err := c.sut.GenerateJWT(expectedUserName)

	userName, err := c.sut.ParseJWT(token)
	assert.Nil(t, err, "expected to complete without error")
	assert.Equal(t, expectedUserName, userName, "resolved user name doesn't match the expected value")
}

// TestTokenUtils_ParseJWT_Invalid_Token tests parsing an expired JWT
func TestTokenUtils_ParseJWT_Invalid_Token(t *testing.T) {
	t.Parallel()
	c := createTokenUtilsContext(t)

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MDEyOTY4MTQsInVzZXIiOiJUZXN0QXV0aG9yIn0.SizHrpPKarCkiJ4taSoUJKNX_GUnT_edSghlFjrgWzg"

	_, err := c.sut.ParseJWT(expiredToken)
	assert.NotNil(t, err, "expired token should lead to error")
	assert.Equal(t, "Token is expired", err.Error(), "incorrect error type")
}

// TestTokenUtils_ParseJWT_Invalid_Signing_Method tests parsing a JWT with incorrect signing method
func TestTokenUtils_ParseJWT_Invalid_Signing_Method(t *testing.T) {
	t.Parallel()
	c := createTokenUtilsContext(t)

	invalidToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3MDIwNzY3ODcsInVzZXIiOiJUZXN0QXV0aG9yIn0.niz4ThjRwsI-2BsoE2F3zLyaAptgt3tWiRzjH_I83q8Uhj0_q6H9N2cFAbkKB_899FQW43pYznocX6oaxM1yZQ"

	_, err := c.sut.ParseJWT(invalidToken)
	assert.NotNil(t, err, "invalid token should lead to error")
	assert.Equal(t, "unexpected signing method: ES256", err.Error(), "incorrect error type")
}

// TestTokenUtils_ParseJWT_Invalid_Claims tests parsing a JWT with missing claims
func TestTokenUtils_ParseJWT_Invalid_Claims(t *testing.T) {
	t.Parallel()
	c := createTokenUtilsContext(t)

	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.LwimMJA3puF3ioGeS-tfczR3370GXBZMIL-bdpu4hOU"

	_, err := c.sut.ParseJWT(invalidToken)
	assert.NotNil(t, err, "invalid token should lead to error")
	assert.Equal(t, "failed to get jwt claims", err.Error(), "incorrect error type")
}
