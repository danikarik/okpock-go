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

func TestCheckUsernameHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *CheckUsernameRequest
		Expected int
	}{
		{
			Name: "NewUser",
			User: nil,
			Request: &CheckUsernameRequest{
				Username: "testuser",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "ExistingUsername",
			User: &testUser{
				Username: "testusernew",
				Email:    "testusernew@example.com",
				Password: "test",
			},
			Request: &CheckUsernameRequest{
				Username: "testusernew",
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

			if tc.User != nil {
				user := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)

				err = srv.env.Auth.SaveNewUser(ctx, user)
				if !assert.NoError(err) {
					return
				}
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			req := newRequest("POST", "/check/username", body, nil, nil)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}
		})
	}
}
