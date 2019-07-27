package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/mux"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler(t *testing.T) {
	assert := assert.New(t)

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	req := newRequest("GET", "/notexists", nil, nil, nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusNotFound, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(err) {
		return
	}
	defer resp.Body.Close()

	var httpErr mux.HTTPError
	err = json.Unmarshal(data, &httpErr)
	if !assert.NoError(err) {
		t.Log(string(data))
		return
	}

	assert.Equal(http.StatusNotFound, httpErr.Code)
	assert.Equal(http.StatusText(http.StatusNotFound), httpErr.Message)
	assert.Equal("404: url=/notexists", httpErr.InternalMessage)
	assert.Nil(httpErr.InternalError)
	assert.False(httpErr.ShowError)
}

func TestMethodNotAllowedHandler(t *testing.T) {
	assert := assert.New(t)

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	req := newRequest("POST", "/", nil, nil, nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusMethodNotAllowed, resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(err) {
		return
	}
	defer resp.Body.Close()

	var httpErr mux.HTTPError
	err = json.Unmarshal(data, &httpErr)
	if !assert.NoError(err) {
		return
	}

	assert.Equal(http.StatusMethodNotAllowed, httpErr.Code)
	assert.Equal(http.StatusText(http.StatusMethodNotAllowed), httpErr.Message)
	assert.Equal("405: url=/, method=POST", httpErr.InternalMessage)
	assert.Nil(httpErr.InternalError)
	assert.False(httpErr.ShowError)
}
