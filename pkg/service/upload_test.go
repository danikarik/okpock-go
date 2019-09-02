package service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/filestore"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUploadHandler(t *testing.T) {
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

			req := authRequest(srv, user, httptest.NewRequest("POST", "/uploads", body))
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

func TestUploadsHandler(t *testing.T) {
	testCases := []struct {
		Name string
		Path string
	}{
		{
			Name: "Image",
			Path: "testdata/gopher.jpg",
		},
		{
			Name: "Text",
			Path: "testdata/test.txt",
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

			hash, err := api.Hash(content)
			if !assert.NoError(err) {
				return
			}

			upload := api.NewUpload(fakeString(), fi.Name(), hash)
			err = srv.env.Logic.SaveNewUpload(ctx, user, upload)
			if !assert.NoError(err) {
				return
			}

			req := authRequest(srv, user, newRequest("GET", "/uploads", nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			var uploads []*api.Upload
			err = unmarshalJSON(resp, &uploads)
			if !assert.NoError(err) {
				return
			}

			if assert.Len(uploads, 1) {
				loaded := uploads[0]
				assert.Equal(upload.ID, loaded.ID)
				assert.Equal(upload.Filename, loaded.Filename)
				assert.Equal(upload.UUID, loaded.UUID)
				assert.Equal(hash, loaded.Hash)
			}
		})
	}
}

func TestUploadHandler(t *testing.T) {
	testCases := []struct {
		Name string
		Path string
	}{
		{
			Name: "Image",
			Path: "testdata/gopher.jpg",
		},
		{
			Name: "Text",
			Path: "testdata/test.txt",
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

			hash, err := api.Hash(content)
			if !assert.NoError(err) {
				return
			}

			upload := api.NewUpload(fakeString(), fi.Name(), hash)
			err = srv.env.Logic.SaveNewUpload(ctx, user, upload)
			if !assert.NoError(err) {
				return
			}

			url := fmt.Sprintf("/uploads/%d", upload.ID)
			req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			var loaded *api.Upload
			err = unmarshalJSON(resp, &loaded)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(upload.ID, loaded.ID)
			assert.Equal(upload.Filename, loaded.Filename)
			assert.Equal(upload.UUID, loaded.UUID)
			assert.Equal(hash, loaded.Hash)
		})
	}
}

func TestUploadFileHandler(t *testing.T) {
	testCases := []struct {
		Name string
		Path string
	}{
		{
			Name: "Image",
			Path: "testdata/gopher.jpg",
		},
		{
			Name: "Text",
			Path: "testdata/test.txt",
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

			hash, err := api.Hash(content)
			if !assert.NoError(err) {
				return
			}

			upload := &api.Upload{
				Filename:    fi.Name(),
				Hash:        hash,
				Body:        content,
				ContentType: http.DetectContentType(content),
				CreatedAt:   time.Now(),
			}

			bucket := srv.env.Config.UploadBucket
			object := &filestore.Object{
				Prefix:      strconv.FormatInt(user.ID, 10),
				Key:         uuid.NewV4().String(),
				Body:        upload.Body,
				ContentType: upload.ContentType,
			}

			err = srv.env.Storage.UploadFile(ctx, bucket, object)
			if !assert.NoError(err) {
				return
			}
			upload.UUID = object.Path()

			err = srv.env.Logic.SaveNewUpload(ctx, user, upload)
			if !assert.NoError(err) {
				return
			}

			url := fmt.Sprintf("/uploads/%d/file", upload.ID)
			req := authRequest(srv, user, newRequest("GET", url, nil, nil, nil))
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)
			resp := rec.Result()

			if !assert.Equal(http.StatusOK, resp.StatusCode) {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if !assert.NoError(err) {
				return
			}
			defer resp.Body.Close()

			assert.Equal(content, data)
		})
	}
}
