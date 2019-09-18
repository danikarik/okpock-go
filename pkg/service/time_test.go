package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsModifiedSince(t *testing.T) {
	testCases := []struct {
		Name            string
		IfModifiedSince string
		LastModified    string
		Expected        bool
	}{
		{
			Name:            "Modified",
			IfModifiedSince: "Tue, 17 Sep 2019 18:35:22 GMT",
			LastModified:    "Tue, 17 Sep 2019 18:39:30 GMT",
			Expected:        false,
		},
		{
			Name:            "Modified",
			IfModifiedSince: "Tue, 17 Sep 2019 19:04:45 GMT",
			LastModified:    "Tue, 17 Sep 2019 19:07:52 GMT",
			Expected:        false,
		},
		{
			Name:            "Modified",
			IfModifiedSince: "Tue, 17 Sep 2019 19:04:45 GMT",
			LastModified:    "Tue, 17 Sep 2019 19:07:52 GMT",
			Expected:        false,
		},
		{
			Name:            "Modified",
			IfModifiedSince: "Tue, 17 Sep 2019 19:04:45 GMT",
			LastModified:    "Tue, 17 Sep 2019 19:05:36 GMT",
			Expected:        false,
		},
		{
			Name:            "NotModified",
			IfModifiedSince: "Tue, 17 Sep 2019 18:35:22 GMT",
			LastModified:    "Tue, 17 Sep 2019 18:35:22 GMT",
			Expected:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("If-Modified-Since", tc.IfModifiedSince)

			lastModified, err := time.Parse(http.TimeFormat, tc.LastModified)
			assert.NoError(err)

			result := notModified(req, lastModified)
			assert.Equal(tc.Expected, result)
		})
	}
}
