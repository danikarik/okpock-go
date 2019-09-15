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

func TestUpdatePassCardHandler(t *testing.T) {
	testCases := []struct {
		Name            string
		UseSerialNumber bool
		Request         *CreatePassCardRequest
	}{
		{
			Name: "UpdateByID",
			Request: &CreatePassCardRequest{
				Structure: &api.PassStructure{
					BackFields: []*api.Field{
						&api.Field{
							Key:   "offer",
							Label: "Any premium dog food",
							Value: "20% off",
						},
					},
				},
			},
		},
		{
			Name:            "UpdateBySerialNumber",
			UseSerialNumber: true,
			Request: &CreatePassCardRequest{
				Structure: &api.PassStructure{
					AuxiliaryFields: []*api.Field{
						&api.Field{
							Key:        "expires",
							Label:      "EXPIRES",
							Value:      "2020-04-24T10:00-05:00",
							IsRelative: true,
							DateStyle:  api.PKDateStyleShort,
						},
					},
				},
			},
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

			project := &api.Project{
				ID:               fakeID(),
				Title:            fakeString(),
				OrganizationName: fakeString(),
				Description:      fakeString(),
				PassType:         api.Coupon,
			}

			err = srv.env.Logic.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			passcard := fakePassCard(project)
			err = srv.env.Logic.SaveNewPassCard(ctx, project, passcard)
			if !assert.NoError(err) {
				return
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			url := fmt.Sprintf("/projects/%d/cards/%d", project.ID, passcard.ID)
			if tc.UseSerialNumber {
				url = fmt.Sprintf("/projects/%d/cards/%s", project.ID, passcard.Data.SerialNumber)
			}
			req := authRequest(srv, user, newRequest("PUT", url, body, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				data := &api.PassCardInfo{}

				err = unmarshalJSON(resp, &data)
				if !assert.NoError(err) {
					return
				}

				assert.True(data.ID > 0)
				assert.Equal(tc.Request.Structure, data.Data.Coupon)
			}
		})
	}
}
