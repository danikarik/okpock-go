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

func TestProjectPassCardHandler(t *testing.T) {
	testCases := []struct {
		Name         string
		SavePassCard bool
		Expected     int
	}{
		{
			Name:         "ExistingPassCard",
			SavePassCard: true,
			Expected:     http.StatusOK,
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

			project := &api.Project{
				ID:               fakeID(),
				Description:      fakeString(),
				OrganizationName: fakeString(),
				PassType:         api.Coupon,
			}

			err = srv.env.Logic.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			passcard := fakePassCard(project)

			if tc.SavePassCard {
				err = srv.env.Logic.SaveNewPassCard(ctx, project, passcard)
				if !assert.NoError(err) {
					return
				}
			}

			url := fmt.Sprintf("/projects/%d/cards/%d", project.ID, passcard.ID)
			req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				data := &api.PassCardInfo{}
				err = unmarshalJSON(resp, &data)
				if !assert.NoError(err) {
					return
				}

				assert.Equal(passcard.ID, data.ID)
				assert.Equal(passcard.Data, data.Data)
			}
		})
	}
}

func TestProjectPassCardBySerialNumberHandler(t *testing.T) {
	testCases := []struct {
		Name         string
		SavePassCard bool
		Expected     int
	}{
		{
			Name:         "ExistingPassCard",
			SavePassCard: true,
			Expected:     http.StatusOK,
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

			project := &api.Project{
				ID:               fakeID(),
				Description:      fakeString(),
				OrganizationName: fakeString(),
				PassType:         api.Coupon,
			}

			err = srv.env.Logic.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			passcard := fakePassCard(project)

			if tc.SavePassCard {
				err = srv.env.Logic.SaveNewPassCard(ctx, project, passcard)
				if !assert.NoError(err) {
					return
				}
			}

			url := fmt.Sprintf("/projects/%d/cards/%s", project.ID, passcard.Data.SerialNumber)
			req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				data := &api.PassCardInfo{}
				err = unmarshalJSON(resp, &data)
				if !assert.NoError(err) {
					return
				}

				assert.Equal(passcard.ID, data.ID)
				assert.Equal(passcard.Data, data.Data)
			}
		})
	}
}
