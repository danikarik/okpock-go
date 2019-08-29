package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/store"
)

// UpdateProjectRequest holds project info to be updated.
type UpdateProjectRequest struct {
	Title            string `json:"title"`
	OrganizationName string `json:"organizationName"`
	Description      string `json:"description"`
}

// IsValid checks whether input is valid or not.
func (r *UpdateProjectRequest) IsValid() error {
	if r.Title == "" {
		return errors.New("titlee is empty")
	}
	if r.OrganizationName == "" {
		return errors.New("organization name is empty")
	}
	if r.Description == "" {
		return errors.New("description is empty")
	}

	return nil
}

// String returns string representation of struct.
func (r *UpdateProjectRequest) String() string {
	return fmt.Sprintf(
		`{"title":"%s","organizationName":"%s","description":"%s"}`,
		r.Title,
		r.OrganizationName,
		r.Description,
	)
}

func (s *Service) updateProjectHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req UpdateProjectRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
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

	err = s.env.Logic.UpdateProject(
		ctx,
		req.Title,
		req.OrganizationName,
		req.Description,
		project,
	)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdateProject", err)
	}

	return sendJSON(w, http.StatusOK, project)
}
