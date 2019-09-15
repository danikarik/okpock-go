package service

import (
	"net/http"

	"github.com/danikarik/mux"
	"github.com/danikarik/okpock/pkg/store"
)

func (s *Service) projectPassCardHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	id, err := s.idFromRequest(r, "id")
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "IDFromRequest", err)
	}

	project, err := s.env.Logic.LoadProject(ctx, user, id)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadProject", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadProject", err)
	}

	cardID, err := s.idFromRequest(r, "cardID")
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "IDFromRequest", err)
	}

	passcard, err := s.env.Logic.LoadPassCard(ctx, project, cardID)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadPassCard", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadPassCard", err)
	}

	return sendJSON(w, http.StatusOK, passcard)
}

func (s *Service) projectPassCardBySerialNumberHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	id, err := s.idFromRequest(r, "id")
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "IDFromRequest", err)
	}

	project, err := s.env.Logic.LoadProject(ctx, user, id)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadProject", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadProject", err)
	}

	passcard, err := s.env.Logic.LoadPassCardBySerialNumber(ctx, project, vars["serialNumber"])
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "LoadPassCardBySerialNumber", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadPassCardBySerialNumber", err)
	}

	return sendJSON(w, http.StatusOK, passcard)
}
