package service

import (
	"errors"
	"fmt"
	"net/http"
)

// UsernameChangeRequest holds new username.
type UsernameChangeRequest struct {
	Username string `json:"username"`
}

// IsValid checks whether input is valid or not.
func (r *UsernameChangeRequest) IsValid() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *UsernameChangeRequest) String() string {
	return fmt.Sprintf(
		`{"username":"%s"}`,
		r.Username,
	)
}

func (s *Service) usernameChangeHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req UsernameChangeRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	// trigger change if emails not match
	if user.Username != req.Username {
		exists, err := s.env.Auth.IsUsernameExists(ctx, req.Username)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "IsUsernameExists", err)
		}
		if exists {
			return sendJSON(w, http.StatusNotAcceptable, M{"username": req.Username})
		}

		exists, err = s.env.Auth.IsEmailExists(ctx, req.Username)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "IsEmailExists", err)
		}
		if exists {
			return sendJSON(w, http.StatusNotAcceptable, M{"email": req.Username})
		}

		err = s.env.Auth.UpdateUsername(ctx, req.Username, user)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "UpdateUsername", err)
		}

		return sendJSON(w, http.StatusOK, M{
			"id":        user.ID,
			"username":  user.Username,
			"updatedAt": user.UpdatedAt,
		})
	}

	return sendJSON(w, http.StatusNotAcceptable, user)
}
