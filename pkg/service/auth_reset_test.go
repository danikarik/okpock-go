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

func TestResetHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name        string
		User        *testUser
		Recover     bool
		NewPassword string
		Token       string
		Expected    int
	}{
		{
			Name: "RecoverFound",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Recover:     true,
			NewPassword: "test_new",
			Expected:    http.StatusAccepted,
		},
		{
			Name: "BadRequest",
			User: &testUser{
				Username: "badrequest",
				Email:    "badrequest@example.com",
				Password: "test",
			},
			Recover:     false,
			NewPassword: "",
			Expected:    http.StatusBadRequest,
		},
		{
			Name: "NotFound",
			User: &testUser{
				Username: "notfound",
				Email:    "notfound@example.com",
				Password: "test",
			},
			Recover:     true,
			NewPassword: "test_new",
			Token:       "qwerty123",
			Expected:    http.StatusNotFound,
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

			user, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			if tc.Recover {
				err = srv.env.Auth.SetRecoveryToken(ctx, user)
				if !assert.NoError(err) {
					return
				}
			}

			token := user.GetRecoveryToken()
			if tc.Token != "" {
				token = tc.Token
			}

			resetRequest := &ResetRequest{
				Token:    token,
				Password: tc.NewPassword,
			}

			body, err := json.Marshal(resetRequest)
			if !assert.NoError(err) {
				return
			}

			req := newRequest("POST", "/reset", body, nil, nil)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusAccepted {
				loaded, err := srv.env.Auth.LoadUserByUsernameOrEmail(ctx, tc.User.Email)
				if !assert.NoError(err) {
					return
				}

				assert.Empty(loaded.GetRecoveryToken())
				assert.True(loaded.CheckPassword(tc.NewPassword))
			}
		})
	}
}
