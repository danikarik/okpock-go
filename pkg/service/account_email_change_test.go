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

func TestEmailChangeHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *EmailChangeRequest
		Expected int
	}{
		{
			Name: "ExistingUser",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Request: &EmailChangeRequest{
				Email: "newtestuser@example.com",
			},
			Expected: http.StatusOK,
		},
		{
			Name: "SameEmail",
			User: &testUser{
				Username: "sameemail",
				Email:    "sameemail@example.com",
				Password: "test",
			},
			Request: &EmailChangeRequest{
				Email: "sameemail@example.com",
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
			Request: &EmailChangeRequest{
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

			user, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
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

			req := newRequest("PUT", "/account/email", body, nil, nil)
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

				assert.Equal(tc.Request.Email, loaded.GetEmailChange())
				assert.NotEmpty(loaded.GetEmailChangeToken())
				assert.NotNil(loaded.EmailChangeSentAt)
			}
		})
	}
}