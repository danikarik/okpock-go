package service

import (
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

func (s *Service) projectPassCardsHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	opts, err := readPagingOptions(r)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadPagingOptions", err)
	}

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

	passcards := &api.PassCardInfoList{}
	searchTerm := r.URL.Query().Get("barcode_message")

	if searchTerm != "" {
		passcards, err = s.env.Logic.LoadPassCardsByBarcodeMessage(ctx, project, searchTerm, opts)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "LoadPassCardsByBarcodeMessage", err)
		}
	} else {
		passcards, err = s.env.Logic.LoadPassCards(ctx, project, opts)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "LoadPassCards", err)
		}
	}

	return sendPaginatedJSON(w, http.StatusOK, passcards.Opts, passcards.Data)
}
