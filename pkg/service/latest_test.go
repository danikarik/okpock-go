package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/stretchr/testify/assert"
)

func TestLatestPass(t *testing.T) {
	ctx := context.Background()

	assert := assert.New(t)

	testCase := struct {
		SerialNumber string
		AuthToken    string
		PassTypeID   string
	}{
		SerialNumber: "9973af9d-9cfa-4d9f-8c6c-32255de8d96b",
		AuthToken:    "secret",
		PassTypeID:   "com.example.pass",
	}

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	err = srv.env.PassKit.InsertPass(
		ctx,
		testCase.SerialNumber,
		testCase.AuthToken,
		testCase.PassTypeID,
	)
	if !assert.NoError(err) {
		return
	}

	path := testCase.SerialNumber + ".pkpass"
	body, err := fakeFile("testdata/" + path)
	if !assert.NoError(err) {
		return
	}

	obj := &filestore.Object{
		Prefix:      "",
		Key:         testCase.SerialNumber,
		Body:        body,
		ContentType: filestore.ApplePkpass,
	}
	err = srv.env.Storage.UploadFile(ctx, srv.env.Config.PassesBucket, obj)
	if !assert.NoError(err) {
		return
	}

	req := newRequest(
		"GET",
		fmt.Sprintf("/v1/passes/%s/%s", testCase.PassTypeID, testCase.SerialNumber),
		nil,
		map[string]string{"Authorization": "ApplePass " + testCase.AuthToken},
		nil,
	)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.NotEmpty(resp.Header.Get("Last-Modified"))
}
