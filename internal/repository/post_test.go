package repository_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wlachs/blog/internal/errortypes"
	"github.com/wlachs/blog/internal/logger"
	"github.com/wlachs/blog/internal/repository"
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

// TestPostRepository_AddPost tests adding a new post to the system
func TestPostRepository_AddPost(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	inputPost := repository.Post{
		URLHandle: "testHandle",
		AuthorID:  0,
	}

	expectedPost := repository.Post{
		URLHandle: inputPost.URLHandle,
	}

	postQuery := regexp.QuoteMeta("INSERT INTO `posts` (`url_handle`,`author_id`,`title`,`summary`,`body`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")
	query := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`url_handle` = ? LIMIT ?")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	c.mockDb.ExpectCommit()
	c.mockDb.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle"}).
			AddRow(expectedPost.ID, expectedPost.URLHandle))

	post, err := c.sut.AddPost(inputPost)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, expectedPost, post, "received post should match the expected one")
}

// TestPostRepository_AddPost_Duplicate_Post tests adding a new post to the system with an already existing URL handle
func TestPostRepository_AddPost_Duplicate_Post(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	inputPost := repository.Post{
		URLHandle: "testHandle",
		AuthorID:  0,
	}

	dbErr := fmt.Errorf("1062")
	expectedError := errortypes.DuplicateElementError{Key: inputPost.URLHandle}

	postQuery := regexp.QuoteMeta("INSERT INTO `posts` (`url_handle`,`author_id`,`title`,`summary`,`body`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnError(dbErr)
	c.mockDb.ExpectRollback()

	post, err := c.sut.AddPost(inputPost)

	assert.Equal(t, repository.Post{}, post, "should not return a post")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_AddPost_Unexpected_Error tests adding a new post to the system while encountering an unexpected error
func TestPostRepository_AddPost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	inputPost := repository.Post{
		URLHandle: "testHandle",
		AuthorID:  0,
	}

	expectedError := fmt.Errorf("unexpected error")

	postQuery := regexp.QuoteMeta("INSERT INTO `posts` (`url_handle`,`author_id`,`title`,`summary`,`body`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	post, err := c.sut.AddPost(inputPost)

	assert.Equal(t, repository.Post{}, post, "should not return a post")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_UpdatePost tests updating a post
func TestPostRepository_UpdatePost(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	title := "newTitle"
	summary := "newSummary"
	body := "newBody"
	inputPost := repository.Post{
		URLHandle: "testHandle",
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}

	expectedPost := repository.Post{
		URLHandle: inputPost.URLHandle,
	}

	postQuery := regexp.QuoteMeta("UPDATE `posts` SET `url_handle`=?,`title`=?,`summary`=?,`body`=?,`updated_at`=? WHERE `posts`.`url_handle` = ?")
	query := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`url_handle` = ? LIMIT ?")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	c.mockDb.ExpectCommit()
	c.mockDb.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle"}).
			AddRow(expectedPost.ID, expectedPost.URLHandle))

	post, err := c.sut.UpdatePost(inputPost)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, expectedPost, post, "received post should match the expected one")
}

// TestPostRepository_UpdatePost_Record_Not_Found tests updating a post non-existing post
func TestPostRepository_UpdatePost_Record_Not_Found(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	title := "newTitle"
	summary := "newSummary"
	body := "newBody"
	inputPost := repository.Post{
		URLHandle: "testHandle",
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}

	postQuery := regexp.QuoteMeta("UPDATE `posts` SET `url_handle`=?,`title`=?,`summary`=?,`body`=?,`updated_at`=? WHERE `posts`.`url_handle` = ?")
	expectedError := errortypes.PostNotFoundError{URLHandle: inputPost.URLHandle}

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	c.mockDb.ExpectCommit()

	post, err := c.sut.UpdatePost(inputPost)

	assert.Equal(t, repository.Post{}, post, "should not return a post")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_UpdatePost_Unexpected_Error tests updating a post with an error
func TestPostRepository_UpdatePost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	title := "newTitle"
	summary := "newSummary"
	body := "newBody"
	inputPost := repository.Post{
		URLHandle: "testHandle",
		Title:     &title,
		Summary:   &summary,
		Body:      &body,
	}

	postQuery := regexp.QuoteMeta("UPDATE `posts` SET `url_handle`=?,`title`=?,`summary`=?,`body`=?,`updated_at`=? WHERE `posts`.`url_handle` = ?")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	post, err := c.sut.UpdatePost(inputPost)

	assert.Equal(t, repository.Post{}, post, "should not return a post")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_DeletePost tests deleting a post
func TestPostRepository_DeletePost(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	urlHandle := "testHandle"

	postQuery := regexp.QuoteMeta("DELETE FROM `posts` WHERE `posts`.`url_handle` = ?")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnResult(sqlmock.NewResult(0, 1))
	c.mockDb.ExpectCommit()

	err := c.sut.DeletePost(urlHandle)

	assert.Nil(t, err, "should complete without error")
}

// TestPostRepository_DeletePost_Record_Not_Found tests deleting a post non-existing post
func TestPostRepository_DeletePost_Record_Not_Found(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	urlHandle := "testHandle"

	postQuery := regexp.QuoteMeta("DELETE FROM `posts` WHERE `posts`.`url_handle` = ?")
	expectedError := errortypes.PostNotFoundError{URLHandle: urlHandle}

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	c.mockDb.ExpectCommit()

	err := c.sut.DeletePost(urlHandle)

	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_DeletePost_Unexpected_Error tests deleting a post with an error
func TestPostRepository_DeletePost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	urlHandle := "testHandle"

	postQuery := regexp.QuoteMeta("DELETE FROM `posts` WHERE `posts`.`url_handle` = ?")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectBegin()
	c.mockDb.ExpectExec(postQuery).WillReturnError(expectedError)
	c.mockDb.ExpectRollback()

	err := c.sut.DeletePost(urlHandle)

	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_GetPost tests retrieving a single post from the database
func TestPostRepository_GetPost(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	expectedPost := repository.Post{
		URLHandle: "testHandle",
	}

	query := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`url_handle` = ? LIMIT ?")

	c.mockDb.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle"}).
			AddRow(expectedPost.ID, expectedPost.URLHandle))

	post, err := c.sut.GetPost(expectedPost.URLHandle)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, expectedPost, post, "received post should match the expected one")
}

// TestPostRepository_GetPost_Record_Not_Found tests retrieving a non-existent post from the database
func TestPostRepository_GetPost_Record_Not_Found(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	expectedPost := repository.Post{
		URLHandle: "testHandle",
	}

	query := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`url_handle` = ? LIMIT ?")
	dbErr := fmt.Errorf("record not found")
	expectedError := errortypes.PostNotFoundError{URLHandle: expectedPost.URLHandle}

	c.mockDb.ExpectQuery(query).WillReturnError(dbErr)

	post, err := c.sut.GetPost(expectedPost.URLHandle)

	assert.Equal(t, repository.Post{}, post, "should not return a post")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_GetPost_Unexpected_Error tests retrieving a single post from the database with an error
func TestPostRepository_GetPost_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `posts` WHERE `posts`.`url_handle` = ? LIMIT ?")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectQuery(query).WillReturnError(expectedError)

	post, err := c.sut.GetPost("test")

	assert.Equal(t, repository.Post{}, post, "should not return a post")
	assert.Equal(t, expectedError, err, "received error should match the expected one")
}

// TestPostRepository_GetPosts tests retrieving every post from the database
func TestPostRepository_GetPosts(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `posts` ORDER BY created_at DESC LIMIT ? OFFSET ?")

	c.mockDb.ExpectQuery(query).
		WithArgs(3, 3).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url_handle"}).
			AddRow(1, "test_1").
			AddRow(2, "test_2"))

	posts, err := c.sut.GetPosts(2, 3)

	assert.Nil(t, err, "should complete without error")
	assert.Equal(t, 2, len(posts), "didn't receive the expected number of posts")
}

// TestPostRepository_GetPosts_Unexpected_Error tests retrieving every post from the database with an error
func TestPostRepository_GetPosts_Unexpected_Error(t *testing.T) {
	t.Parallel()
	c := createPostRepositoryContext(t)

	query := regexp.QuoteMeta("SELECT * FROM `posts` ORDER BY created_at DESC")
	expectedError := fmt.Errorf("unexpected error")

	c.mockDb.ExpectQuery(query).WillReturnError(expectedError)

	posts, err := c.sut.GetPosts(1, 1)

	assert.Equal(t, expectedError, err, "error should match expected value")
	assert.Equal(t, 0, len(posts), "shouldn't receive any posts")
}
