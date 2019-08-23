package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
)

// RegisterRequest holds auth credentials to register.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// IsValid checks whether input is valid or not.
func (r *RegisterRequest) IsValid() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	if r.Email == "" {
		return errors.New("email is empty")
	}
	if r.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *RegisterRequest) String() string {
	return fmt.Sprintf(
		`{"username":"%s","email":"%s","password":"%s"}`,
		r.Username,
		r.Email,
		r.Password,
	)
}

func (s *Service) registerHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req RegisterRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	exists, err := s.env.Auth.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsUsernameExists", err)
	}
	if exists {
		return sendJSON(w, http.StatusNotAcceptable, M{"username": req.Username})
	}

	exists, err = s.env.Auth.IsEmailExists(ctx, req.Email)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsEmailExists", err)
	}
	if exists {
		return sendJSON(w, http.StatusNotAcceptable, M{"email": req.Email})
	}

	hash, err := secure.NewPassword(req.Password)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "NewPassword", err)
	}

	user := api.NewUser(req.Username, req.Email, hash, map[string]interface{}{})

	err = s.env.Auth.SaveNewUser(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SaveNewUser", err)
	}

	err = s.env.Auth.SetConfirmationToken(ctx, api.SignUpConfirmation, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SetConfirmationToken", err)
	}

	message, err := s.confirmMessage(user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "ConfirmMessage", err)
	}

	sentAt, err := s.env.Mailer.SendMail(ctx, message)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SendMail", err)
	}

	return sendJSON(w, http.StatusCreated, M{
		"email":     user.Email,
		"messageId": message.ID,
		"sentAt":    sentAt,
	})
}
