package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	assert := assert.New(t)

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	user, err := api.NewUser(fakeUsername(), fakeUsername(), "test", nil)
	if !assert.NoError(err) {
		return
	}

	err = srv.env.Auth.SaveNewUser(context.Background(), user)
	if !assert.NoError(err) {
		return
	}

	ucl := NewClaims().WithUser(user).WithCSRFToken(newCSRFToken())
	tokenString, _ := ucl.MarshalJWT()
	tokenCookie := srv.tokenCookie(tokenString)

	req := newRequest("DELETE", "/logout", nil, nil, nil)
	req.AddCookie(tokenCookie)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusForbidden, resp.StatusCode)

	req.Header.Add(csrfHeader, ucl.CSRFToken)
	rec = httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp = rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCORSMiddleware(t *testing.T) {
	assert := assert.New(t)

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}
	srv.env.Config.Stage = "development"

	req := newRequest("OPTIONS", "/login", nil, nil, nil)
	req.Header.Set("Origin", "http://localhost:8080")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("POST", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal("http://localhost:8080", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal("Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
}
