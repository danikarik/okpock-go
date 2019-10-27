package service

import "net/http"

func (s *Service) authCheckHandler(w http.ResponseWriter, r *http.Request) error {
	var (
		ucl    = NewClaims()
		status = http.StatusOK
	)

	err := s.getClaims(r, ucl)
	if err != nil {
		status = http.StatusUnauthorized
	}

	return sendJSON(w, status, M{"ping": "pong"})
}
