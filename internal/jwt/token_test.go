package jwt_test

import (
	"github.com/wlchs/blog/internal/jwt"
	"github.com/wlchs/blog/internal/logger"
	"testing"
)

// jwtTestContext contains objects relevant for testing the TokenUtils.
type jwtTestContext struct {
	sut jwt.TokenUtils
}

// createJWTServiceContext creates the context for testing the TokenUtils and reduces code duplication.
func createJWTServiceContext(t *testing.T) *jwtTestContext {
	t.Helper()

	sut := jwt.CreateJWTUtils(logger.CreateLogger())

	return &jwtTestContext{sut}
}
