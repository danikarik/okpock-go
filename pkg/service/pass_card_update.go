package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danikarik/mux"
	"github.com/danikarik/okpock/pkg/store"
)

func (s *Service) updatePassCardHandler(w http.ResponseWriter, r *http.Request) error {
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

	pushToken, err := s.env.PassKit.FindPushToken(ctx, passcard.Data.SerialNumber)
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "FindPushToken", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "FindPushToken", err)
	}

	var req CreatePassCardRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	newPasscard, err := s.newProjectPassCard(&req, project)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "NewProjectPassCard", err)
	}

	err = newPasscard.IsValid()
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	err = s.env.Logic.UpdatePassCard(ctx, newPasscard.Data, passcard)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdatePassCard", err)
	}

	err = s.env.PassKit.UpdatePass(ctx, passcard.Data.SerialNumber)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdatePass", err)
	}

	notificator, err := s.getNotificator(project.PassType)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "GetNotificator", err)
	}

	err = notificator.Push(ctx, pushToken)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Push", err)
	}

	return sendJSON(w, http.StatusOK, passcard)
}

func (s *Service) updatePassCardBySerialNumberHandler(w http.ResponseWriter, r *http.Request) error {
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
		return s.httpError(w, r, http.StatusNotFound, "LoadPassCard", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LoadPassCard", err)
	}

	pushToken, err := s.env.PassKit.FindPushToken(ctx, vars["serialNumber"])
	if err == store.ErrNotFound {
		return s.httpError(w, r, http.StatusNotFound, "FindPushToken", err)
	}
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "FindPushToken", err)
	}

	var req CreatePassCardRequest
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	newPasscard, err := s.newProjectPassCard(&req, project)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "NewProjectPassCard", err)
	}

	err = newPasscard.IsValid()
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	err = s.env.Logic.UpdatePassCard(ctx, newPasscard.Data, passcard)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdatePassCard", err)
	}

	err = s.env.PassKit.UpdatePass(ctx, passcard.Data.SerialNumber)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdatePass", err)
	}

	notificator, err := s.getNotificator(project.PassType)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "GetNotificator", err)
	}

	err = notificator.Push(ctx, pushToken)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Push", err)
	}

	return sendJSON(w, http.StatusOK, passcard)
}
