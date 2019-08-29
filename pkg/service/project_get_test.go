package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUserProjectHandler(t *testing.T) {
	testCases := []struct {
		Name        string
		SaveProject bool
		Expected    int
	}{
		{
			Name:        "ExistingProject",
			SaveProject: true,
			Expected:    http.StatusOK,
		},
		{
			Name:     "NotFound",
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

			user := api.NewUser(fakeUsername(), fakeEmail(), fakePassword(), nil)
			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := api.NewProject(fakeString(),
				fakeString(),
				fakeString(),
				api.Coupon)

			if tc.SaveProject {
				err = srv.env.Logic.SaveNewProject(ctx, user, project)
				if !assert.NoError(err) {
					return
				}
			} else {
				project.ID = 1 << 32
			}

			url := fmt.Sprintf("/projects/%d", project.ID)
			req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				data := &api.Project{}
				err = unmarshalJSON(resp, &data)
				if !assert.NoError(err) {
					return
				}

				assert.Equal(project.ID, data.ID)
				assert.Equal(project.Title, data.Title)
				assert.Equal(project.OrganizationName, data.OrganizationName)
				assert.Equal(project.Description, data.Description)
				assert.Equal(project.PassType, data.PassType)
			}
		})
	}
}
