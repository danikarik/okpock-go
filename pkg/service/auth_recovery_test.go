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

func TestRecoverHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *RecoverRequest
		Expected int
	}{
		{
			Name: "ExistingUser",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Request: &RecoverRequest{
				Email: "testuser@example.com",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "NotFound",
			User: nil,
			Request: &RecoverRequest{
				Email: "testuser@example.com",
			},
			Expected: http.StatusNotFound,
		},
		{
			Name: "EmptyEmail",
			User: nil,
			Request: &RecoverRequest{
				Email: "",
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

			req := newRequest("POST", "/recover", body, nil, nil)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				loaded, err := srv.env.Auth.LoadUserByUsernameOrEmail(ctx, tc.Request.Email)
				if !assert.NoError(err) {
					return
				}

				assert.NotEmpty(loaded.RecoveryToken)
				assert.NotNil(loaded.RecoverySentAt)
			}
		})
	}
}
