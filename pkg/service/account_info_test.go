package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestAccountHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Expected int
	}{
		{
			Name: "ExistingUser",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Expected: http.StatusOK,
		},
		{
			Name:     "NotFound",
			User:     nil,
			Expected: http.StatusUnauthorized,
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

			req := newRequest("GET", "/account/info", nil, nil, nil)
			rec := httptest.NewRecorder()

			if tc.User != nil {
				user := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)

				err = srv.env.Auth.SaveNewUser(ctx, user)
				if !assert.NoError(err) {
					return
				}

				ucl := NewClaims().WithUser(user).WithCSRFToken(newCSRFToken())
				tokenString, _ := ucl.MarshalJWT()
				tokenCookie := srv.tokenCookie(tokenString)

				req.Header.Set(csrfHeader, ucl.CSRFToken)
				req.AddCookie(tokenCookie)
			}

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				data, err := ioutil.ReadAll(resp.Body)
				if !assert.NoError(err) {
					return
				}
				defer resp.Body.Close()

				var loaded api.User
				err = json.Unmarshal(data, &loaded)
				if !assert.NoError(err) {
					return
				}
				assert.Equal(tc.User.Email, loaded.Email)
			}
		})
	}
}
