package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/store/sequel"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestSerialNumbers(t *testing.T) {
	ctx := context.Background()

	assert := assert.New(t)

	deviceID := uuid.NewV4().String()
	passTypeID := "com.example.pass"
	now := time.Now()

	testCases := []struct {
		SerialNumber string
		DeviceID     string
		AuthToken    string
		PassTypeID   string
		UpdatedAt    time.Time
	}{
		{
			SerialNumber: uuid.NewV4().String(),
			DeviceID:     deviceID,
			AuthToken:    uuid.NewV4().String(),
			PassTypeID:   passTypeID,
			UpdatedAt:    now.Add(3 * time.Second),
		},
		{
			SerialNumber: uuid.NewV4().String(),
			DeviceID:     deviceID,
			AuthToken:    uuid.NewV4().String(),
			PassTypeID:   passTypeID,
			UpdatedAt:    now.Add(5 * time.Second),
		},
	}

	srv, err := initService(t)
	if !assert.NoError(err) {
		return
	}

	for _, tc := range testCases {
		err = srv.env.PassKit.InsertPass(
			ctx,
			tc.SerialNumber,
			tc.AuthToken,
			tc.PassTypeID,
		)
		if !assert.NoError(err) {
			return
		}

		err = srv.env.PassKit.InsertRegistration(
			ctx,
			tc.DeviceID,
			uuid.NewV4().String(),
			tc.SerialNumber,
			tc.PassTypeID,
		)
		if !assert.NoError(err) {
			return
		}
	}

	values := url.Values{}
	values.Set("passesUpdatedSince", now.Add(-4*time.Second).Format(sequel.TimeFormat))
	req := newRequest(
		"GET",
		fmt.Sprintf("/v1/devices/%s/registrations/%s", deviceID, passTypeID),
		nil,
		nil,
		values,
	)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)
	resp := rec.Result()

	assert.Equal(http.StatusOK, resp.StatusCode)
}
