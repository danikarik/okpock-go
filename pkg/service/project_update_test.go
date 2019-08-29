package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProjectHandler(t *testing.T) {
	testCases := []struct {
		Name     string
		Request  *UpdateProjectRequest
		Expected int
	}{
		{
			Name: "Coupon",
			Request: &UpdateProjectRequest{
				Title:            "Saturday Deal",
				OrganizationName: "Okpock Child",
				Description:      "Free Auction",
			},
			Expected: http.StatusOK,
		},
		{
			Name:     "Invalid",
			Request:  &UpdateProjectRequest{},
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

			user := api.NewUser(fakeUsername(), fakeEmail(), fakePassword(), nil)
			err = srv.env.Auth.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			project := api.NewProject(fakeString(), fakeString(), fakeString(), api.Coupon)
			err = srv.env.Logic.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			url := fmt.Sprintf("/projects/%d", project.ID)
			req := authRequest(srv, user, newRequest("PUT", url, body, nil, nil))
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

				assert.True(data.ID > 0)
				assert.Equal(tc.Request.Title, data.Title)
				assert.Equal(tc.Request.OrganizationName, data.OrganizationName)
				assert.Equal(tc.Request.Description, data.Description)
			}
		})
	}
}
