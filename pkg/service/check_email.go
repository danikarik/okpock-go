package service

import (
	"errors"
	"fmt"
	"net/http"
)

// CheckEmailRequest holds email to be checked.
type CheckEmailRequest struct {
	Email string `json:"email"`
}

// IsValid checks whether input is valid or not.
func (r *CheckEmailRequest) IsValid() error {
	if r.Email == "" {
		return errors.New("email is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *CheckEmailRequest) String() string {
	return fmt.Sprintf(
		`{"email":"%s"}`,
		r.Email,
	)
}

func (s *Service) checkEmailHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req CheckEmailRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	exists, err := s.env.Auth.IsEmailExists(ctx, req.Email)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsEmailExists", err)
	}

	code := http.StatusOK
	if exists {
		code = http.StatusNotAcceptable
	}

	return sendJSON(w, code, M{"email": req.Email})
}
