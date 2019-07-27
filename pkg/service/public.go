package service

import "net/http"

func (s *Service) okHandler(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	return nil
}
