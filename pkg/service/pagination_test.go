package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestPaginatedJSON(t *testing.T) {
	testCases := []struct {
		Name string
		Opts *api.PagingOptions
		List []interface{}
	}{
		{
			Name: "Projects",
			Opts: &api.PagingOptions{
				Cursor: 0,
				Limit:  3,
				Next:   4,
			},
			List: []interface{}{
				api.NewProject(fakeString(), fakeString(), fakeString(), api.Coupon),
				api.NewProject(fakeString(), fakeString(), fakeString(), api.BoardingPass),
				api.NewProject(fakeString(), fakeString(), fakeString(), api.EventTicket),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert := assert.New(t)

			rec := httptest.NewRecorder()
			err := sendPaginatedJSON(rec, 200, tc.Opts, tc.List)
			if !assert.NoError(err) {
				return
			}

			resp := rec.Result()
			data, err := ioutil.ReadAll(resp.Body)
			if !assert.NoError(err) {
				return
			}
			defer resp.Body.Close()

			var payload M
			err = json.Unmarshal(data, &payload)
			if !assert.NoError(err) {
				return
			}

			token, ok := payload["token"].(string)
			if !assert.True(ok) {
				return
			}
			if !assert.NotEmpty(token) {
				return
			}

			list, ok := payload["data"].([]interface{})
			if !assert.True(ok) {
				return
			}
			if !assert.NotNil(list) {
				return
			}
			if !assert.Len(list, len(tc.List)) {
				return
			}

			var limit uint64 = 10
			query := url.Values{}
			query.Add(pageTokenQuery, token)
			query.Add(pageLimitQuery, fmt.Sprintf("%d", limit))

			req := httptest.NewRequest("GET", "/", nil)
			req.URL.RawQuery = query.Encode()

			opts, err := readPagingOptions(req)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(tc.Opts.Next, opts.Cursor)
			assert.Equal(limit, opts.Limit)
		})
	}
}
