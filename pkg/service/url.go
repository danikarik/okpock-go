package service

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/danikarik/mux"
	"github.com/danikarik/okpock/pkg/api"
)

var (
	// ErrMissingQueryParam raised when one of required query parameters is missing.
	ErrMissingQueryParam = errors.New("url: missing query parameter")
	// ErrInvalidID raises when cannot parse id from request.
	ErrInvalidID = errors.New("id: invalid query parameter")
)

func checkQueryParams(r *http.Request, params ...string) (map[string]string, error) {
	vars := mux.Vars(r)
	for _, param := range params {
		if v, ok := vars[param]; !ok || v == "" {
			return vars, ErrMissingQueryParam
		}
	}
	return vars, nil
}

func (s *Service) hostURL() string {
	if s.env.Config.Debug {
		if host, ok := os.LookupEnv("HOST_URL"); ok {
			return host
		}
		return "http://localhost:" + s.env.Config.Port
	}
	if s.env.Config.IsDevelopment() {
		return "https://api-dev.okpock.com"
	}
	if s.env.Config.IsProduction() {
		return "https://api.okpock.com"
	}
	return ""
}

func (s *Service) appURL(path string) string {
	if s.env.Config.Debug {
		return "http://localhost:3000" + path
	}
	if s.env.Config.IsDevelopment() {
		return "https://console-dev.okpock.com" + path
	}
	if s.env.Config.IsProduction() {
		return "https://console.okpock.com" + path
	}
	return ""
}

func (s *Service) confirmationURL(u *api.User, c api.Confirmation) (string, error) {
	link, err := url.Parse(s.hostURL() + "/verify")
	if err != nil {
		return "", err
	}

	values := url.Values{}
	switch c {
	case api.SignUpConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.ConfirmationToken)
		values.Add("redirect_url", s.appURL(""))
		break
	case api.InviteConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.ConfirmationToken)
		values.Add("redirect_url", s.appURL("/reset"))
		break
	case api.RecoveryConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.RecoveryToken)
		values.Add("redirect_url", s.appURL("/reset"))
		break
	case api.EmailChangeConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.EmailChangeToken)
		values.Add("redirect_url", s.appURL(""))
		break
	}

	link.RawQuery = values.Encode()
	return link.String(), nil
}

func (s *Service) redirect(w http.ResponseWriter, r *http.Request, url string) error {
	http.Redirect(w, r, url, http.StatusMovedPermanently)
	return nil
}

func (s *Service) redirectError(w http.ResponseWriter, r *http.Request, msg string, err error) error {
	url, uerr := url.Parse(s.appURL("/error"))
	if uerr != nil {
		return uerr
	}
	v := url.Query()
	v.Add("message", msg)
	if err != nil && !s.env.Config.IsProduction() {
		v.Add("error", err.Error())
	}
	url.RawQuery = v.Encode()
	http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
	return nil
}

func (s *Service) idFromRequest(r *http.Request, key string) (int64, error) {
	vars := mux.Vars(r)

	v, ok := vars[key]
	if !ok {
		return -1, ErrInvalidID
	}

	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1, ErrInvalidID
	}

	return id, nil
}
