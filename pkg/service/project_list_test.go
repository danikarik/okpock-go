package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUserProjectsHandler(t *testing.T) {
	testCases := []struct {
		Name          string
		ProjectNumber int
	}{
		{Name: "EmptyProjectList"},
		{Name: "NotEmptyProjectList", ProjectNumber: 10},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			srv, err := initService(t)
			if !assert.NoError(err) {
				return
			}

			user := api.NewUser(fakeUsername(), fakeEmail(), fakePassword(), nil)
			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			for i := 0; i < tc.ProjectNumber; i++ {
				project := api.NewProject(fakeString(),
					fakeString(),
					fakeString(),
					api.Coupon)

				err = srv.env.Logic.SaveNewProject(ctx, user, project)
				if !assert.NoError(err) {
					return
				}
			}

			req := authRequest(srv, user, newRequest("GET", "/projects/", nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			data := []*api.Project{}
			err = unmarshalJSON(resp, &data)
			if !assert.NoError(err) {
				return
			}

			assert.Len(data, tc.ProjectNumber)
		})
	}
}
