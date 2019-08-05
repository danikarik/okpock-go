package service

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

const recoveryTokenTTL = 1 * time.Hour

// ErrExpiredToken raises when token TTL exceeded.
var ErrExpiredToken = errors.New("token: allowed time is expired")

func tokenExpired(t time.Time) error {
	if time.Since(t) > recoveryTokenTTL {
		return ErrExpiredToken
	}
	return nil
}

func (s *Service) verifyHandler(w http.ResponseWriter, r *http.Request) error {
	vars, err := checkQueryParams(r)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "CheckQueryParams", err)
	}

	confirm := api.Confirmation(vars["type"])
	switch confirm {
	case api.SignUpConfirmation:
		return s.verifyBySignUpConfirmationToken(vars, w, r)
	case api.InviteConfirmation:
		return s.verifyByInviteConfirmationToken(vars, w, r)
	case api.RecoveryConfirmation:
		return s.verifyByRecoveryToken(vars, w, r)
	case api.EmailChangeConfirmation:
		return s.verifyByEmailChangeToken(vars, w, r)
	}

	return s.httpError(w, r, http.StatusBadRequest, "Confirmation", api.ErrUnknownConfirmation)
}

func (s *Service) verifyBySignUpConfirmationToken(vars map[string]string, w http.ResponseWriter, r *http.Request) error {
	var (
		ctx         = r.Context()
		token       = vars["token"]
		redirectURL = vars["redirect_url"]
	)

	user, err := s.env.Auth.LoadUserByConfirmationToken(ctx, token)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadUserByConfirmationToken", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByConfirmationToken", err)
	}

	err = s.env.Auth.ConfirmUser(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "ConfirmUser", err)
	}

	return s.redirect(w, r, redirectURL)
}

func (s *Service) verifyByInviteConfirmationToken(vars map[string]string, w http.ResponseWriter, r *http.Request) error {
	var (
		token       = vars["token"]
		confirm     = vars["type"]
		redirectURL = vars["redirect_url"]
	)

	url, err := url.Parse(redirectURL)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Parse", err)
	}

	v := url.Query()
	v.Add("type", confirm)
	v.Add("token", token)
	url.RawQuery = v.Encode()

	return s.redirect(w, r, url.String())
}

func (s *Service) verifyByRecoveryToken(vars map[string]string, w http.ResponseWriter, r *http.Request) error {
	var (
		ctx         = r.Context()
		token       = vars["token"]
		confirm     = vars["type"]
		redirectURL = vars["redirect_url"]
	)

	user, err := s.env.Auth.LoadUserByRecoveryToken(ctx, token)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByRecoveryToken", err)
	}

	err = tokenExpired(*user.RecoverySentAt)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "TokenExpired", err)
	}

	url, err := url.Parse(redirectURL)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Parse", err)
	}

	v := url.Query()
	v.Add("type", confirm)
	v.Add("token", token)
	url.RawQuery = v.Encode()

	return s.redirect(w, r, url.String())
}

func (s *Service) verifyByEmailChangeToken(vars map[string]string, w http.ResponseWriter, r *http.Request) error {
	var (
		ctx         = r.Context()
		token       = vars["token"]
		redirectURL = vars["redirect_url"]
	)

	user, err := s.env.Auth.LoadUserByEmailChangeToken(ctx, token)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadUserByEmailChangeToken", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByEmailChangeToken", err)
	}

	err = s.env.Auth.ConfirmEmailChange(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "ConfirmEmailChange", err)
	}

	err = s.clearCookies(w)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "ClearCookies", err)
	}

	return s.redirect(w, r, redirectURL)
}
