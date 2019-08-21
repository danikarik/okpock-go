package service

import (
	"net/http"
)

func (s *Service) accountInfoHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := userFromContext(r.Context())
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}
	return sendJSON(w, http.StatusOK, user)
}
