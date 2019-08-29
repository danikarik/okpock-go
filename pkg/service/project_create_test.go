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

func TestCreateProjectHandler(t *testing.T) {
	testCases := []struct {
		Name       string
		SaveBefore bool
		Request    *CheckProjectRequest
		Expected   int
	}{
		{
			Name:       "NewProject",
			SaveBefore: false,
			Request: &CheckProjectRequest{
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         "coupon",
			},
			Expected: http.StatusCreated,
		},
		{
			Name:       "ExistingProject",
			SaveBefore: true,
			Request: &CheckProjectRequest{
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         "boardingPass",
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

			user := api.NewUser(fakeUsername(), fakeEmail(), fakePassword(), nil)
			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			if tc.SaveBefore {
				project := api.NewProject(
					tc.Request.Title,
					tc.Request.OrganizationName,
					tc.Request.Description,
					api.PassType(tc.Request.PassType))

				err = srv.env.Logic.SaveNewProject(ctx, user, project)
				if !assert.NoError(err) {
					return
				}
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			req := authRequest(srv, user, newRequest("POST", "/projects/", body, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusCreated {
				data := M{}
				err = unmarshalJSON(resp, &data)
				if !assert.NoError(err) {
					return
				}

				id := int64(data["id"].(float64))
				loaded, err := srv.env.Logic.LoadProject(ctx, user, id)
				if !assert.NoError(err) {
					return
				}

				assert.Equal(tc.Request.Title, loaded.Title)
				assert.Equal(tc.Request.OrganizationName, loaded.OrganizationName)
				assert.Equal(tc.Request.Description, loaded.Description)
				assert.Equal(api.PassType(tc.Request.PassType), loaded.PassType)
			}
		})
	}
}
