package service

import (
	"errors"
	"fmt"
	"net/http"
)

// CheckUsernameRequest holds username to be checked.
type CheckUsernameRequest struct {
	Username string `json:"username"`
}

// IsValid checks whether input is valid or not.
func (r *CheckUsernameRequest) IsValid() error {
	if r.Username == "" {
		return errors.New("username is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *CheckUsernameRequest) String() string {
	return fmt.Sprintf(
		`{"username":"%s"}`,
		r.Username,
	)
}

func (s *Service) checkUsernameHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req CheckUsernameRequest
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

	return sendJSON(w, http.StatusOK, M{"username": req.Username})
}
