package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
)

// CreateControllerContext creates a Context for mocking HTTP requests.
func CreateControllerContext() (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = &http.Request{Header: make(http.Header)}
	return ctx, recorder
}

// MockJsonPost method to mock HTTP request bodies.
// Code from https://stackoverflow.com/questions/57733801/how-to-set-mock-gin-context-for-bindjson
func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST" // or PUT
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
