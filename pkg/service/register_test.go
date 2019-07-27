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

func TestRegisterDevice(t *testing.T) {
	assert := assert.New(t)

	testCase := struct {
		SerialNumber string
		DeviceID     string
		AuthToken    string
		PassTypeID   string
	}{
		SerialNumber: uuid.NewV4().String(),
		DeviceID:     uuid.NewV4().String(),
		AuthToken:    uuid.NewV4().String(),
		PassTypeID:   "com.example.pass",
	}

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	err = srv.env.PassKit.InsertPass(
		context.Background(),
		testCase.SerialNumber,
		testCase.AuthToken,
		testCase.PassTypeID,
	)
	if !assert.NoError(err) {
		return
	}

	testStatusCodes := []int{http.StatusCreated, http.StatusOK}
	for _, code := range testStatusCodes {
		req := newRequest(
			"POST",
			fmt.Sprintf("/v1/devices/%s/registrations/%s/%s", testCase.DeviceID, testCase.PassTypeID, testCase.SerialNumber),
			[]byte(`{"pushToken":"test-token"}`),
			map[string]string{"Authorization": "ApplePass " + testCase.AuthToken},
			nil,
		)
		rec := httptest.NewRecorder()

		srv.ServeHTTP(rec, req)
		resp := rec.Result()

		assert.Equal(code, resp.StatusCode)
	}
}
