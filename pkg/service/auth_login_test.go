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

func TestLoginHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *LoginRequest
		Confirm  bool
		Expected int
	}{
		{
			Name: "Confirmed",
			User: &testUser{
				Username: "confirmed",
				Email:    "confirmed@example.com",
				Password: "test",
			},
			Request: &LoginRequest{
				Username: "confirmed",
				Password: "test",
			},
			Confirm:  true,
			Expected: http.StatusOK,
		},
		{
			Name: "NotConfirmed",
			User: &testUser{
				Username: "notconfirmed",
				Email:    "notconfirmed@example.com",
				Password: "test",
			},
			Request: &LoginRequest{
				Username: "notconfirmed",
				Password: "test",
			},
			Expected: http.StatusLocked,
		},
		{
			Name: "NoCredentials",
			User: &testUser{
				Username: "nocredentials",
				Email:    "nocredentials@example.com",
				Password: "test",
			},
			Request: &LoginRequest{
				Username: "",
				Password: "",
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "WrongPassword",
			User: &testUser{
				Username: "wrongpassword",
				Email:    "wrongpassword@example.com",
				Password: "test",
			},
			Request: &LoginRequest{
				Username: "wrongpassword",
				Password: "test2",
			},
			Confirm:  true,
			Expected: http.StatusForbidden,
		},
		{
			Name: "NotFound",
			User: nil,
			Request: &LoginRequest{
				Username: "notfound",
				Password: "test",
			},
			Expected: http.StatusNotFound,
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
				user, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
				if !assert.NoError(err) {
					return
				}

				err = srv.env.Auth.SaveNewUser(ctx, user)
				if !assert.NoError(err) {
					return
				}

				if tc.Confirm {
					err = srv.env.Auth.ConfirmUser(ctx, user)
					if !assert.NoError(err) {
						return
					}
				}
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			req := newRequest("POST", "/login", body, nil, nil)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			assert.Equal(tc.Expected, resp.StatusCode)
		})
	}
}
