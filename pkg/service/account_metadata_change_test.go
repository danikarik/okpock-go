package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestMetaDataChangeHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *MetaDataChangeRequest
		Expected int
	}{
		{
			Name: "ExistingUser",
			User: &testUser{
				Username: "testuser",
				Email:    "testuser@example.com",
				Password: "test",
			},
			Request: &MetaDataChangeRequest{
				Data: api.JSONMap{
					"lastSeen": time.Now(),
				},
			},
			Expected: http.StatusOK,
		},
		{
			Name: "EmptyData",
			User: &testUser{
				Username: "emptyemail",
				Email:    "emptyemail@example.com",
				Password: "test",
			},
			Request: &MetaDataChangeRequest{
				Data: api.JSONMap{},
			},
			Expected: http.StatusBadRequest,
		},
		{
			Name: "NilData",
			User: &testUser{
				Username: "nildata",
				Email:    "nildata@example.com",
				Password: "test",
			},
			Request: &MetaDataChangeRequest{
				Data: nil,
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

			req := newRequest("PUT", "/account/metadata", body, nil, nil)
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

				foundAll := true
				for k := range tc.Request.Data {
					if _, ok := loaded.UserMetaData[k]; !ok {
						foundAll = false
					}
				}
				assert.True(foundAll)
			}
		})
	}
}
