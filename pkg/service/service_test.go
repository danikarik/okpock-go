package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func fakeUsername() string {
	return uuid.NewV4().String()
}

func fakeEmail() string {
	return uuid.NewV4().String() + "@example.com"
}

func fakeFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func initService(t *testing.T) (*Service, error) {
	_, err := env.NewLookup(
		"TEST_DATABASE_URL",
		"TEST_PASSES_BUCKET",
		"TEST_SERVER_SECRET",
	)
	if err != nil {
		t.Skip(err)
	}

	e := env.NewMock()

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return New("test-version", e, logger), nil
}

func newRequest(method, url string, body []byte, headers map[string]string, values url.Values) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	req.URL.RawQuery = values.Encode()
	req.Header = newHeader(headers)
	return req
}

func newHeader(h map[string]string) http.Header {
	header := make(http.Header)
	if h != nil && len(h) > 0 {
		for k, v := range h {
			header.Set(k, v)
		}
	}
	return header
}
