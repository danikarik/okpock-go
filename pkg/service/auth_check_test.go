package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/require"
)

func TestAuthCheckHandler(t *testing.T) {
	testCases := []struct {
		Name      string
		SetCookie bool
		Expected  int
	}{
		{
			Name:      "Authorized",
			SetCookie: true,
			Expected:  http.StatusOK,
		},
		{
			Name:      "Unauthorized",
			SetCookie: false,
			Expected:  http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			require := require.New(t)

			srv, err := initService(t)
			require.NoError(err)

			user := api.NewUser(
				fakeUsername(),
				fakeEmail(),
				fakePassword(),
				nil,
			)
			err = srv.env.Auth.SaveNewUser(ctx, user)
			require.NoError(err)

			ucl := NewClaims().WithUser(user).WithCSRFToken(newCSRFToken())
			tokenString, _ := ucl.MarshalJWT()
			tokenCookie := srv.tokenCookie(tokenString)

			req := newRequest("GET", "/ping", nil, nil, nil)
			if tc.SetCookie {
				req.AddCookie(tokenCookie)
			}
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			require.Equal(tc.Expected, resp.StatusCode)
		})
	}
}
