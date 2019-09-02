package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUsernameChangeHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *UsernameChangeRequest
		Expected int
	}{
		{
			Name: "ExistingUser",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Request: &UsernameChangeRequest{
				Username: "newtestuser",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "SameUsername",
			User: &testUser{
				Username: "sameemail",
				Email:    "sameemail@example.com",
				Password: "test",
			},
			Request: &UsernameChangeRequest{
				Username: "sameemail",
			},
			Expected: http.StatusNotAcceptable,
		},
		{
			Name: "EmptyEmail",
			User: &testUser{
				Username: "emptyemail",
				Email:    "emptyemail@example.com",
				Password: "test",
			},
			Request: &UsernameChangeRequest{
				Username: "",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "Duplicated",
			User: &testUser{
				Username: "duplicated",
				Email:    "duplicated@example.com",
				Password: "test",
			},
			Request: &UsernameChangeRequest{
				Username: "another",
			},
			Expected: http.StatusNotAcceptable,
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

			another := api.NewUser("another", "another@example.com", "test", nil)

			err = srv.env.Auth.SaveNewUser(ctx, another)
			if !assert.NoError(err) {
				return
			}

			user := &api.User{
				ID:           fakeID(),
				Username:     tc.User.Username,
				Email:        tc.User.Email,
				PasswordHash: tc.User.Password,
			}

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

			req := newRequest("PUT", "/account/username", body, nil, nil)
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
				if !assert.NoError(err) {
					return
				}

				assert.Equal(tc.Request.Username, loaded.Username)
			}
		})
	}
}
