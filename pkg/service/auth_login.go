package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

// LoginRequest holds auth credentials.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// IsValid checks whether input is valid or not.
func (r *LoginRequest) IsValid() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *LoginRequest) String() string {
	return fmt.Sprintf(
		`{"username":"%s","password":"%s"}`,
		r.Username,
		r.Password,
	)
}

func withUserClaims(w http.ResponseWriter, u *api.User) error {
	ucl := NewClaims().
		WithUser(u).
		WithCSRFToken(newCSRFToken())
	err := setClaimsCookie(w, ucl)
	if err != nil {
		return err
	}
	w.Header().Set(csrfHeader, ucl.CSRFToken)
	return nil
}

func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req LoginRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := s.env.Auth.LoadUserByUsernameOrEmail(ctx, req.Username)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadUserByUsernameOrEmail", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByUsernameOrEmail", err)
	}

	err = s.env.Auth.Authenticate(ctx, req.Password, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Authenticate", err)
	}

	err = withUserClaims(w, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "WithUserClaims", err)
	}

	return nil
}

func (s *Service) logoutHandler(w http.ResponseWriter, r *http.Request) error {
	return clearCookies(w)
}
