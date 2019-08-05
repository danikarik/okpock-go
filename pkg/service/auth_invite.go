package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
)

// InviteRequest holds referer and email to invite.
type InviteRequest struct {
	Email string `json:"email"`
}

// IsValid checks whether input is valid or not.
func (r *InviteRequest) IsValid() error {
	if r.Email == "" {
		return errors.New("email is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *InviteRequest) String() string {
	return fmt.Sprintf(
		`{"email":"%s"}`,
		r.Email,
	)
}

func (s *Service) inviteHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req InviteRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	authUser, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	exists, err := s.env.Auth.IsUsernameExists(ctx, req.Email)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsUsernameExists", err)
	}
	if exists {
		return sendJSON(w, http.StatusNotAcceptable, M{"username": req.Email})
	}

	exists, err = s.env.Auth.IsEmailExists(ctx, req.Email)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsEmailExists", err)
	}
	if exists {
		return sendJSON(w, http.StatusNotAcceptable, M{"email": req.Email})
	}

	// use email as email itself and username.
	user, err := api.NewUser(req.Email, req.Email, secure.Token(), map[string]interface{}{})
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "NewUser", err)
	}

	err = s.env.Auth.SaveNewUser(ctx, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SaveNewUser", err)
	}

	err = s.env.Auth.SetConfirmationToken(ctx, api.InviteConfirmation, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SetConfirmationToken", err)
	}

	err = s.env.Auth.UpdateAppMetaData(ctx, M{
		metaReferer:               authUser.Email,
		metaSuggestChangeUsername: true,
	}, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdateAppMetaData", err)
	}

	message, err := s.inviteMessage(authUser, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "InviteMessage", err)
	}

	sentAt, err := s.env.Mailer.SendMail(ctx, message)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SendMail", err)
	}

	return sendJSON(w, http.StatusCreated, M{
		"id":        user.ID,
		"email":     user.Email,
		"messageId": message.ID,
		"sentAt":    sentAt,
	})
}
