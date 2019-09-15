package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/secure"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func fakeID() int64 {
	return rand.Int63n(10000)
}

func fakeString() string {
	return uuid.NewV4().String()
}

func fakeUsername() string {
	return uuid.NewV4().String()
}

func fakeEmail() string {
	return uuid.NewV4().String() + "@example.com"
}

func fakePassword() string {
	raw := "test"
	hash, err := secure.NewPassword(raw)
	if err != nil {
		return raw
	}
	return hash
}

func fakeFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func fakePassCard(project *api.Project) *api.PassCardInfo {
	now := time.Now()
	data := &api.PassCard{
		Description:         project.Description,
		FormatVersion:       1,
		OrganizationName:    project.OrganizationName,
		PassTypeID:          "pass.com.okpock.test",
		SerialNumber:        uuid.NewV4().String(),
		TeamID:              fakeString(),
		WebServiceURL:       "http://localhost:5000",
		AuthenticationToken: secure.Token(),
	}
	return &api.PassCardInfo{
		ID:        fakeID(),
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func initService(t *testing.T) (*Service, error) {
	_, err := env.NewLookup(
		"TEST_UPLOAD_BUCKET",
		"TEST_PASSES_BUCKET",
		"TEST_SERVER_SECRET",
	)
	if err != nil {
		t.Skip(err)
	}

	e, err := env.NewMock()
	if err != nil {
		t.Fatal(err)
	}

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
func authRequest(srv *Service, u *api.User, req *http.Request) *http.Request {
	ucl := NewClaims().WithUser(u).WithCSRFToken(newCSRFToken())
	tokenString, _ := ucl.MarshalJWT()
	tokenCookie := srv.tokenCookie(tokenString)
	req.Header.Set(csrfHeader, ucl.CSRFToken)
	req.AddCookie(tokenCookie)
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

func unmarshalJSON(r *http.Response, v interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}
