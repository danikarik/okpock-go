package service

import (
	"errors"
	"fmt"
	"net/http"
)

// EmailChangeRequest holds new email.
type EmailChangeRequest struct {
	Email string `json:"email"`
}

// IsValid checks whether input is valid or not.
func (r *EmailChangeRequest) IsValid() error {
	if r.Email == "" {
		return errors.New("Email is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *EmailChangeRequest) String() string {
	return fmt.Sprintf(
		`{"email":"%s"}`,
		r.Email,
	)
}

func (s *Service) emailChangeHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req EmailChangeRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	// trigger change if emails not match
	if user.Email != req.Email {
		err = s.env.Auth.SetEmailChangeToken(ctx, req.Email, user)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SetEmailChangeToken", err)
		}

		message, err := s.emailChangeMessage(user)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "EmailChangeMessage", err)
		}

		sentAt, err := s.env.Mailer.SendMail(ctx, message)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SendMail", err)
		}

		return sendJSON(w, http.StatusOK, M{
			"id":        user.ID,
			"email":     user.Email,
			"messageId": message.ID,
			"sentAt":    sentAt,
		})
	}

	return sendJSON(w, http.StatusNotAcceptable, user)
}
