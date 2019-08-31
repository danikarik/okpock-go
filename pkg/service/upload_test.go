package service

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	testCases := []struct {
		Name       string
		SaveBefore bool
		Path       string
		Expected   int
	}{
		{
			Name:     "NewUpload",
			Path:     "testdata/gopher.jpg",
			Expected: http.StatusCreated,
		},
		{
			Name:       "Duplicated",
			SaveBefore: true,
			Path:       "testdata/gopher.jpg",
			Expected:   http.StatusNotAcceptable,
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

			if tc.SaveBefore {
				hash, err := api.Hash(content)
				if !assert.NoError(err) {
					return
				}

				upload := api.NewUpload(fakeString(), fi.Name(), hash)
				err = srv.env.Logic.SaveNewUpload(ctx, user, upload)
				if !assert.NoError(err) {
					return
				}
			}

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", fi.Name())
			if !assert.NoError(err) {
				return
			}

			_, err = part.Write(content)
			if !assert.NoError(err) {
				return
			}

			err = writer.Close()
			if !assert.NoError(err) {
				return
			}

			req := authRequest(srv, user, httptest.NewRequest("POST", "/upload", body))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(tc.Expected, resp.StatusCode) {
				return
			}

			if resp.StatusCode == http.StatusOK {
				var upload *api.Upload
				err = unmarshalJSON(resp, &upload)
				if !assert.NoError(err) {
					return
				}
			}
		})
	}
}
