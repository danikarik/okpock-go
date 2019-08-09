package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

// ResetRequest holds user password and recovery token.
type ResetRequest struct {
	Type     string `json:"type"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

// IsValid checks whether input is valid or not.
func (r *ResetRequest) IsValid() error {
	if r.Type == "" {
		return errors.New("type is empty")
	}
	if r.Token == "" {
		return errors.New("token is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *ResetRequest) String() string {
	return fmt.Sprintf(
		`{"type":"%s","token":"%s","password":"%s"}`,
		r.Type,
		r.Token,
		r.Password,
	)
}

func (s *Service) resetHandler(w http.ResponseWriter, r *http.Request) error {
	var req ResetRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	confirm := api.Confirmation(req.Type)
	switch confirm {
	case api.InviteConfirmation:
		return s.resetByConfirmationToken(req, w, r)
	case api.RecoveryConfirmation:
		return s.resetByRecoveryToken(req, w, r)
	}

	return s.httpError(w, r, http.StatusBadRequest, "Confirmation", api.ErrUnknownConfirmation)

}

func (s *Service) resetByConfirmationToken(req ResetRequest, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	user, err := s.env.Auth.LoadUserByConfirmationToken(ctx, req.Token)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadUserByConfirmationToken", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByConfirmationToken", err)
	}

	err = s.env.Auth.UpdatePassword(ctx, req.Password, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdatePassword", err)
	}

	err = s.env.Auth.ConfirmUser(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "ConfirmUser", err)
	}

	return sendJSON(w, http.StatusAccepted, M{"confirmedAt": user.ConfirmedAt})
}

func (s *Service) resetByRecoveryToken(req ResetRequest, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	user, err := s.env.Auth.LoadUserByRecoveryToken(ctx, req.Token)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadUserByRecoveryToken", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByRecoveryToken", err)
	}

	err = s.env.Auth.UpdatePassword(ctx, req.Password, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdatePassword", err)
	}

	err = s.env.Auth.RecoverUser(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "RecoverUser", err)
	}

	return sendJSON(w, http.StatusAccepted, M{"updatedAt": user.UpdatedAt})
}
