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
		ID       int64
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
				ID:       10,
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
				ID:       11,
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
				ID:       12,
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
				ID:       13,
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
			another.ID = 999

			err = srv.env.Auth.SaveNewUser(ctx, another)
			if !assert.NoError(err) {
				return
			}

			user := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			user.ID = tc.User.ID

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
