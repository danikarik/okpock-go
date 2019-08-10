package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	assert := assert.New(t)

	user, err := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)
	if !assert.NoError(err) {
		return
	}

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	err = srv.env.Auth.SaveNewUser(context.Background(), user)
	if !assert.NoError(err) {
		return
	}

	raw := fmt.Sprintf(`{"username":"%s","password":"%s"}`, user.Username, "test")
	req := newRequest("POST", "/login", []byte(raw), nil, nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)

	ucl := NewClaims().WithUser(user).WithCSRFToken(resp.Header.Get("X-XSRF-TOKEN"))
	tokenString, _ := ucl.MarshalJWT()
	tokenCookie := srv.tokenCookie(tokenString)

	req = newRequest("DELETE", "/logout", nil, nil, nil)
	req.AddCookie(tokenCookie)
	req.Header.Set("X-XSRF-TOKEN", ucl.CSRFToken)
	rec = httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp = rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)

	req = newRequest("GET", "/account", nil, nil, nil)
	rec = httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp = rec.Result()

	assert.Equal(http.StatusUnauthorized, resp.StatusCode)
}
