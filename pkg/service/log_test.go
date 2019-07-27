package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorLogs(t *testing.T) {
	assert := assert.New(t)

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	req := newRequest(
		"POST",
		"/v1/log",
		[]byte(`{"logs":["test1","test2"]}`),
		nil,
		nil,
	)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
}
