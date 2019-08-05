package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestVerifyHandler(t *testing.T) {
	type testUser struct {
		Username string
		Email    string
		Password string
	}

	testCases := []struct {
		Name     string
		User     *testUser
		Confirm  api.Confirmation
		Expected int
	}{
		{
			Name: "SignUp",
			User: &testUser{
				Username: "signup",
				Email:    "signup@example.com",
				Password: "signup",
			},
			Confirm:  api.SignUpConfirmation,
			Expected: http.StatusMovedPermanently,
		},
		{
			Name: "Invite",
			User: &testUser{
				Username: "invite",
				Email:    "invite@example.com",
				Password: "invite",
			},
			Confirm:  api.InviteConfirmation,
			Expected: http.StatusMovedPermanently,
		},
		{
			Name: "Recovery",
			User: &testUser{
				Username: "recovery",
				Email:    "recovery@example.com",
				Password: "recovery",
			},
			Confirm:  api.RecoveryConfirmation,
			Expected: http.StatusMovedPermanently,
		},
		{
			Name: "EmailChange",
			User: &testUser{
				Username: "emailchange",
				Email:    "emailchange@example.com",
				Password: "emailchange",
			},
			Confirm:  api.EmailChangeConfirmation,
			Expected: http.StatusMovedPermanently,
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

			app := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token := r.URL.Query().Get("token")
				assert.NotEmpty(token)
			}))
			defer app.Close()

			user, err := api.NewUser(tc.User.Username, tc.User.Email, tc.User.Password, nil)
			if !assert.NoError(err) {
				return
			}

			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			var token string
			switch tc.Confirm {
			case api.SignUpConfirmation, api.InviteConfirmation:
				err = srv.env.Auth.SetConfirmationToken(ctx, tc.Confirm, user)
				if !assert.NoError(err) {
					return
				}
				token = user.GetConfirmationToken()
				break
			case api.RecoveryConfirmation:
				err = srv.env.Auth.SetRecoveryToken(ctx, user)
				if !assert.NoError(err) {
					return
				}
				token = user.GetRecoveryToken()
				break
			case api.EmailChangeConfirmation:
				err = srv.env.Auth.SetEmailChangeToken(ctx, "new@example.com", user)
				if !assert.NoError(err) {
					return
				}
				token = user.GetEmailChangeToken()
				break
			}

			url, err := url.Parse("/verify")
			if !assert.NoError(err) {
				return
			}
			v := url.Query()
			v.Add("type", string(tc.Confirm))
			v.Add("token", token)
			v.Add("redirect_url", app.URL)

			req := newRequest("GET", "/verify", nil, nil, v)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusMovedPermanently {
				var loaded *api.User
				{
					if tc.Confirm == api.EmailChangeConfirmation {
						loaded, err = srv.env.Auth.LoadUserByUsernameOrEmail(ctx, "new@example.com")
						if !assert.NoError(err) {
							return
						}
					} else {
						loaded, err = srv.env.Auth.LoadUserByUsernameOrEmail(ctx, tc.User.Email)
						if !assert.NoError(err) {
							return
						}
					}
				}

				switch tc.Confirm {
				case api.SignUpConfirmation:
					assert.Empty(loaded.GetConfirmationToken())
					assert.NotNil(loaded.ConfirmationSentAt)
					break
				case api.InviteConfirmation:
					// should be empty after reset handler
					assert.NotEmpty(loaded.GetConfirmationToken())
					assert.NotNil(loaded.InvitedAt)
					break
				case api.RecoveryConfirmation:
					// should be empty after reset handler
					assert.NotEmpty(loaded.GetRecoveryToken())
					assert.NotNil(loaded.RecoverySentAt)
					break
				case api.EmailChangeConfirmation:
					assert.Empty(loaded.GetEmailChangeToken())
					assert.NotNil(loaded.EmailChangeSentAt)
					assert.Equal("new@example.com", loaded.Email)
					break
				}
			}

		})
	}
}
