package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/secure"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestPasswordChangeHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *PasswordChangeRequest
		Expected int
	}{
		{
			Name: "ExistingUser",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Request: &PasswordChangeRequest{
				Password: "newpass",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "SamePassword",
			User: &testUser{
				Username: "samepassword",
				Email:    "samepassword@example.com",
				Password: "test",
			},
			Request: &PasswordChangeRequest{
				Password: "test",
			},
			Expected: http.StatusNotAcceptable,
		},
		{
			Name: "EmptyPassword",
			User: &testUser{
				Username: "emptypassword",
				Email:    "emptypassword@example.com",
				Password: "test",
			},
			Request: &PasswordChangeRequest{
				Password: "",
			},
			Expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			srv, err := initService(t)
			if !assert.NoError(err) {
				return
			}

			hash, err := secure.NewPassword(tc.User.Password)
			if !assert.NoError(err) {
				return
			}

			user := api.NewUser(tc.User.Username, tc.User.Email, hash, nil)

			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			ucl := NewClaims().WithUser(user).WithCSRFToken(newCSRFToken())
			tokenString, _ := ucl.MarshalJWT()
			tokenCookie := srv.tokenCookie(tokenString)

			req := newRequest("PUT", "/account/password", body, nil, nil)
			req.Header.Set(csrfHeader, ucl.CSRFToken)
			req.AddCookie(tokenCookie)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				loaded, err := srv.env.Auth.LoadUserByUsernameOrEmail(ctx, tc.User.Email)
				assert.NoError(err)
				assert.True(loaded.CheckPassword(tc.Request.Password))
			}
		})
	}
}
