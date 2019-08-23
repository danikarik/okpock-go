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

func TestInviteHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Request  *InviteRequest
		Expected int
	}{
		{
			Name: "NewUser",
			User: nil,
			Request: &InviteRequest{
				Email: "testuser@example.com",
			},
			Expected: http.StatusCreated,
		},
		{
			Name: "ExistedEmail",
			User: &testUser{
				Username: "testuserold",
				Email:    "testuserold@example.com",
				Password: "test",
			},
			Request: &InviteRequest{
				Email: "testuserold@example.com",
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

			referer := api.NewUser(fakeUsername(), fakeEmail(), "test", nil)
			referer.ID = 700

			err = srv.env.Auth.SaveNewUser(ctx, referer)
			if !assert.NoError(err) {
				return
			}

			if tc.User != nil {
				user := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
				user.ID = 707

				err = srv.env.Auth.SaveNewUser(ctx, user)
				if !assert.NoError(err) {
					return
				}
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			ucl := NewClaims().WithUser(referer).WithCSRFToken(newCSRFToken())
			tokenString, _ := ucl.MarshalJWT()
			tokenCookie := srv.tokenCookie(tokenString)

			req := newRequest("POST", "/invite", body, nil, nil)
			req.Header.Set(csrfHeader, ucl.CSRFToken)
			req.AddCookie(tokenCookie)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusCreated {
				loaded, err := srv.env.Auth.LoadUserByUsernameOrEmail(ctx, tc.Request.Email)
				if !assert.NoError(err) {
					return
				}

				assert.Equal(loaded.Username, tc.Request.Email)
				assert.Equal(loaded.Email, tc.Request.Email)
				assert.False(loaded.IsConfirmed())
				assert.NotEmpty(loaded.GetConfirmationToken())
				assert.NotNil(loaded.InvitedAt)

				ref, ok := loaded.AppMetaData[metaReferer]
				assert.True(ok)
				assert.Equal(referer.Email, ref)
			}
		})
	}
}
