package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wlchs/blog/internal/logger"
	"github.com/wlchs/blog/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

// postTestContext contains objects relevant for testing the PostRepository.
type postTestContext struct {
	mockDb sqlmock.Sqlmock
	sut    repository.PostRepository
}

// createPostRepositoryContext creates the context for testing the PostRepository and reduces code duplication.
func createPostRepositoryContext(t *testing.T) *postTestContext {
	t.Helper()

	db, mock, _ := sqlmock.New()
	gormDb, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))

	sut := repository.CreatePostRepository(logger.CreateLogger(), repository.CreateRepository(gormDb))
	return &postTestContext{mock, sut}
}

// TestPostRepository_GetPosts tests retrieving every post from the database
func TestPostRepository_GetPosts(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `posts` ORDER BY created_at DESC")
	c.mockDb.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle"}).
			AddRow(1, "test_1").
			AddRow(2, "test_2"))

	posts, err := c.sut.GetPosts()

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, 2, len(posts), "didn't receive the expected number of posts")
}
