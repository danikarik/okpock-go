package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestUnregisterDevice(t *testing.T) {
	ctx := context.Background()

	assert := assert.New(t)

	testCase := struct {
		SerialNumber string
		DeviceID     string
		AuthToken    string
		PushToken    string
		PassTypeID   string
	}{
		SerialNumber: uuid.NewV4().String(),
		DeviceID:     uuid.NewV4().String(),
		AuthToken:    uuid.NewV4().String(),
		PushToken:    uuid.NewV4().String(),
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

	err = srv.env.PassKit.InsertRegistration(
		ctx,
		testCase.DeviceID,
		testCase.PushToken,
		testCase.SerialNumber,
		testCase.PassTypeID,
	)
	if !assert.NoError(err) {
		return
	}

	req := newRequest(
		"DELETE",
		fmt.Sprintf("/v1/devices/%s/registrations/%s/%s", testCase.DeviceID, testCase.PassTypeID, testCase.SerialNumber),
		nil,
		map[string]string{"Authorization": "ApplePass " + testCase.AuthToken},
		nil,
	)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
}
