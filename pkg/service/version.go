package service

import "net/http"

func (s *Service) versionHandler(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, M{"version": s.version})
}
