package service

import "net/http"

func (s *Service) userProjectsHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	projects, err := s.env.Logic.LoadProjects(ctx, user, nil)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadProjects", err)
	}

	return sendJSON(w, http.StatusOK, projects)
}
