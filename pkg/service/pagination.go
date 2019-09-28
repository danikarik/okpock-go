package service

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danikarik/okpock/pkg/api"
)

const (
	pageTokenQuery = "page_token"
	pageLimitQuery = "page_limit"
)

func readPagingOptions(r *http.Request) (*api.PagingOptions, error) {
	opts := api.NewPagingOptions(0, 0)

	if token := r.URL.Query().Get(pageTokenQuery); token != "" {
		data, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			return nil, err
		}

		o := new(api.PagingOptions)
		err = json.Unmarshal(data, &o)
		if err != nil {
			return nil, err
		}

		opts.Cursor = o.Next
	}

	if limit := r.URL.Query().Get(pageLimitQuery); limit != "" {
		lim, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return nil, err
		}

		opts.Limit = lim
	}

	return opts, nil
}

func sendPaginatedJSON(w http.ResponseWriter, code int, opts *api.PagingOptions, data interface{}) error {
	tdata, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	token := ""
	if opts.HasNext() {
		token = base64.URLEncoding.EncodeToString(tdata)
	}

	return sendJSON(w, code, M{
		"token": token,
		"data":  data,
	})
}
