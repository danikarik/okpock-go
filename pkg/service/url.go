package service

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/danikarik/mux"

	"github.com/danikarik/okpock/pkg/api"
)

// ErrMissingQueryParam raised when one of required query parameters is missing.
var ErrMissingQueryParam = errors.New("url: missing query parameter")

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
		return "http://localhost:" + s.env.Config.Port
	}
	if s.env.Config.IsDevelopment() {
		return "https://api-dev.okpock.com"
	}
	if s.env.Config.IsDevelopment() {
		return "https://api.okpock.com"
	}
	return ""
}

func (s *Service) appURL() string {
	if s.env.Config.Debug {
		return "http://localhost:3000"
	}
	if s.env.Config.IsDevelopment() {
		return "https://app-dev.okpock.com"
	}
	if s.env.Config.IsDevelopment() {
		return "https://app.okpock.com"
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
		values.Add("token", u.GetConfirmationToken())
		values.Add("redirect_url", s.appURL())
		break
	case api.InviteConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.GetConfirmationToken())
		values.Add("redirect_url", s.appURL())
		break
	case api.RecoveryConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.GetRecoveryToken())
		values.Add("redirect_url", s.appURL())
		break
	case api.EmailChangeConfirmation:
		values.Add("type", string(c))
		values.Add("token", u.GetEmailChangeToken())
		values.Add("redirect_url", s.appURL())
		break
	}

	link.RawQuery = values.Encode()
	return link.String(), nil
}

func (s *Service) redirect(w http.ResponseWriter, r *http.Request, url string) error {
	http.Redirect(w, r, url, http.StatusMovedPermanently)
	return nil
}
