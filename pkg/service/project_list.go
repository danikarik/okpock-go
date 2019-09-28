package service

import "net/http"

func (s *Service) userProjectsHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	opts, err := readPagingOptions(r)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadPagingOptions", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	projects, err := s.env.Logic.LoadProjects(ctx, user, opts)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadProjects", err)
	}

	return sendPaginatedJSON(w, http.StatusOK, projects.Opts, projects.Data)
}
