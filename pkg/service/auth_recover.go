package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/store"
)

// RecoverRequest holds recover email.
type RecoverRequest struct {
	Email string `json:"email"`
}

// IsValid checks whether input is valid or not.
func (r *RecoverRequest) IsValid() error {
	if r.Email == "" {
		return errors.New("email is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *RecoverRequest) String() string {
	return fmt.Sprintf(`{"email":"%s"}`, r.Email)
}

func (s *Service) recoverHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req RecoverRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := s.env.Auth.LoadUserByUsernameOrEmail(ctx, req.Email)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadUserByUsernameOrEmail", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadUserByUsernameOrEmail", err)
	}

	err = s.env.Auth.SetRecoveryToken(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SetRecoveryToken", err)
	}

	message, err := s.recoverMessage(user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "RecoverMessage", err)
	}

	sentAt, err := s.env.Mailer.SendMail(ctx, message)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SendMail", err)
	}

	return sendJSON(w, http.StatusOK, M{"messageId": message.ID, "sentAt": sentAt})
}
