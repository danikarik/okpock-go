package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUploadProjectImage(t *testing.T) {
	testCases := []struct {
		Name    string
		Request *UploadImageRequest
		Path    string
	}{
		{
			Name:    "Background",
			Request: &UploadImageRequest{Type: backgroundImage},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Footer",
			Request: &UploadImageRequest{Type: footerImage},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Icon",
			Request: &UploadImageRequest{Type: iconImage},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Logo",
			Request: &UploadImageRequest{Type: logoImage},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Strip",
			Request: &UploadImageRequest{Type: stripImage},
			Path:    "testdata/gopher.jpg",
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

			file, err := os.Open(tc.Path)
			if !assert.NoError(err) {
				return
			}
			defer file.Close()

			content, err := ioutil.ReadAll(file)
			if !assert.NoError(err) {
				return
			}

			fi, err := file.Stat()
			if !assert.NoError(err) {
				return
			}

			hash, err := api.Hash(content)
			if !assert.NoError(err) {
				return
			}

			upload := api.NewUpload(fakeString(), fi.Name(), hash)
			err = srv.env.Logic.SaveNewUpload(ctx, user, upload)
			if !assert.NoError(err) {
				return
			}

			tc.Request.UUID = upload.UUID
			body, err := json.Marshal(tc.Request)
			if !assert.NoError(err) {
				return
			}

			url := fmt.Sprintf("/projects/%d/upload", project.ID)
			req := authRequest(srv, user, newRequest("POST", url, body, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			data := &api.Project{}

			err = unmarshalJSON(resp, &data)
			if !assert.NoError(err) {
				return
			}

			switch tc.Request.Type {
			case backgroundImage:
				assert.Equal(tc.Request.UUID, data.BackgroundImage)
				break
			case footerImage:
				assert.Equal(tc.Request.UUID, data.FooterImage)
				break
			case iconImage:
				assert.Equal(tc.Request.UUID, data.IconImage)
				break
			case logoImage:
				assert.Equal(tc.Request.UUID, data.LogoImage)
				break
			case stripImage:
				assert.Equal(tc.Request.UUID, data.StripImage)
				break
			}
		})
	}
}
