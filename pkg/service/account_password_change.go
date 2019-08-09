package service

import (
	"errors"
	"fmt"
	"net/http"
)

// PasswordChangeRequest holds new password.
type PasswordChangeRequest struct {
	Password string `json:"password"`
}

// IsValid checks whether input is valid or not.
func (r *PasswordChangeRequest) IsValid() error {
	if r.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *PasswordChangeRequest) String() string {
	return fmt.Sprintf(
		`{"password":"%s"}`,
		r.Password,
	)
}

func (s *Service) passwordChangeHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req PasswordChangeRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	// trigger change if password has diff
	if !user.CheckPassword(req.Password) {
		err = s.env.Auth.UpdatePassword(ctx, req.Password, user)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "UpdatePassword", err)
		}

		return sendJSON(w, http.StatusOK, M{
			"updatedAt": user.UpdatedAt,
		})
	}

	return sendJSON(w, http.StatusNotAcceptable, user)
}
