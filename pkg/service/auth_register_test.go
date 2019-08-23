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

func TestRegisterHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *RegisterRequest
		Expected int
	}{
		{
			Name: "NewUser",
			User: nil,
			Request: &RegisterRequest{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "ExistedUsername",
			User: &testUser{
				Username: "testuser",
				Email:    "testusernew@example.com",
				Password: "test",
			},
			Request: &RegisterRequest{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Expected: http.StatusNotAcceptable,
		},
		{
			Name: "ExistedEmail",
			User: &testUser{
				Username: "testusernew",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Request: &RegisterRequest{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
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

			req := newRequest("POST", "/register", body, nil, nil)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusCreated {
				loaded, err := srv.env.Auth.LoadUserByUsernameOrEmail(ctx, tc.Request.Username)
				if !assert.NoError(err) {
					return
				}

				assert.Equal(loaded.Username, tc.Request.Username)
				assert.Equal(loaded.Email, tc.Request.Email)
				assert.True(loaded.CheckPassword(tc.Request.Password))
				assert.False(loaded.IsConfirmed())
				assert.NotEmpty(loaded.GetConfirmationToken())
				assert.NotNil(loaded.ConfirmationSentAt)
			}
		})
	}
}
