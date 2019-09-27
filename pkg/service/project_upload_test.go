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
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/stretchr/testify/assert"
)

func TestUploadProjectImage(t *testing.T) {
	testCases := []struct {
		Name    string
		Request *UploadImageRequest
		Path    string
	}{
		{
			Name:    "Background1x",
			Request: &UploadImageRequest{Type: backgroundImage, Size: api.ImageSize1x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Background2x",
			Request: &UploadImageRequest{Type: backgroundImage, Size: api.ImageSize2x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Background3x",
			Request: &UploadImageRequest{Type: backgroundImage, Size: api.ImageSize3x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Footer1x",
			Request: &UploadImageRequest{Type: footerImage, Size: api.ImageSize1x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Footer2x",
			Request: &UploadImageRequest{Type: footerImage, Size: api.ImageSize2x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Footer3x",
			Request: &UploadImageRequest{Type: footerImage, Size: api.ImageSize3x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Icon1x",
			Request: &UploadImageRequest{Type: iconImage, Size: api.ImageSize1x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Icon2x",
			Request: &UploadImageRequest{Type: iconImage, Size: api.ImageSize2x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Icon3x",
			Request: &UploadImageRequest{Type: iconImage, Size: api.ImageSize3x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Logo1x",
			Request: &UploadImageRequest{Type: logoImage, Size: api.ImageSize1x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Logo2x",
			Request: &UploadImageRequest{Type: logoImage, Size: api.ImageSize2x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Logo3x",
			Request: &UploadImageRequest{Type: logoImage, Size: api.ImageSize3x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Strip1x",
			Request: &UploadImageRequest{Type: stripImage, Size: api.ImageSize1x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Strip2x",
			Request: &UploadImageRequest{Type: stripImage, Size: api.ImageSize2x},
			Path:    "testdata/gopher.jpg",
		},
		{
			Name:    "Strip3x",
			Request: &UploadImageRequest{Type: stripImage, Size: api.ImageSize3x},
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

			hash, err := secure.Hash(content)
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
				switch tc.Request.Size {
				case api.ImageSize1x:
					assert.Equal(tc.Request.UUID, data.BackgroundImage)
				case api.ImageSize2x:
					assert.Equal(tc.Request.UUID, data.BackgroundImage2x)
				case api.ImageSize3x:
					assert.Equal(tc.Request.UUID, data.BackgroundImage3x)
				}
			case footerImage:
				switch tc.Request.Size {
				case api.ImageSize1x:
					assert.Equal(tc.Request.UUID, data.FooterImage)
				case api.ImageSize2x:
					assert.Equal(tc.Request.UUID, data.FooterImage2x)
				case api.ImageSize3x:
					assert.Equal(tc.Request.UUID, data.FooterImage3x)
				}
			case iconImage:
				switch tc.Request.Size {
				case api.ImageSize1x:
					assert.Equal(tc.Request.UUID, data.IconImage)
				case api.ImageSize2x:
					assert.Equal(tc.Request.UUID, data.IconImage2x)
				case api.ImageSize3x:
					assert.Equal(tc.Request.UUID, data.IconImage3x)
				}
			case logoImage:
				switch tc.Request.Size {
				case api.ImageSize1x:
					assert.Equal(tc.Request.UUID, data.LogoImage)
				case api.ImageSize2x:
					assert.Equal(tc.Request.UUID, data.LogoImage2x)
				case api.ImageSize3x:
					assert.Equal(tc.Request.UUID, data.LogoImage3x)
				}
			case stripImage:
				switch tc.Request.Size {
				case api.ImageSize1x:
					assert.Equal(tc.Request.UUID, data.StripImage)
				case api.ImageSize2x:
					assert.Equal(tc.Request.UUID, data.StripImage2x)
				case api.ImageSize3x:
					assert.Equal(tc.Request.UUID, data.StripImage3x)
				}
			}
		})
	}
}
