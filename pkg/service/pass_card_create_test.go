package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/stretchr/testify/assert"
)

func TestCreatePassCardHandler(t *testing.T) {
	testCases := []struct {
		Name     string
		PassType api.PassType
		Request  *CreatePassCardRequest
		Expected int
	}{
		{
			Name:     "Coupon",
			PassType: api.Coupon,
			Request: &CreatePassCardRequest{
				Barcodes: []*api.Barcode{
					&api.Barcode{
						Message:         "123456789",
						Format:          api.PKBarcodeFormatPDF417,
						MessageEncoding: "iso-8859-1",
					},
				},
				Locations: []*api.Location{
					&api.Location{
						Longitude: -122.3748889,
						Latitude:  37.6189722,
					},
					&api.Location{
						Longitude: -122.03118,
						Latitude:  37.33182,
					},
				},
				LogoText:        "Paw Planet",
				ForegroundColor: "rgb(255, 255, 255)",
				BackgroundColor: "rgb(206, 140, 53)",
				Structure: &api.PassStructure{
					PrimaryFields: []*api.Field{
						&api.Field{
							Key:   "offer",
							Label: "Any premium dog food",
							Value: "20% off",
						},
					},
					AuxiliaryFields: []*api.Field{
						&api.Field{
							Key:        "expires",
							Label:      "EXPIRES",
							Value:      "2013-04-24T10:00-05:00",
							IsRelative: true,
							DateStyle:  api.PKDateStyleShort,
						},
					},
				},
			},
			Expected: http.StatusCreated,
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

			project := api.NewProject(fakeString(), fakeString(), fakeString(), api.PassType(tc.PassType))
			err = srv.env.Logic.SaveNewProject(ctx, user, project)
			if !assert.NoError(err) {
				return
			}

			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			rawurl := fmt.Sprintf("/projects/%d/cards", project.ID)
			req := authRequest(srv, user, newRequest("POST", rawurl, body, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusCreated {
				var data = M{}
				err = unmarshalJSON(resp, &data)
				if !assert.NoError(err) {
					return
				}

				downloadURL, err := url.Parse(data["url"].(string))
				if !assert.NoError(err) {
					return
				}

				req = newRequest("GET", downloadURL.Path, nil, nil, nil)
				rec = httptest.NewRecorder()

				srv.ServeHTTP(rec, req)
				resp = rec.Result()

				content, err := ioutil.ReadAll(resp.Body)
				if !assert.NoError(err) {
					return
				}
				defer resp.Body.Close()

				assert.True(len(content) > 0)
				assert.Equal(resp.Header.Get("Content-Type"), filestore.ApplePkpass)
			}
		})
	}
}
