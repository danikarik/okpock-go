package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
)

// CreateProjectRequest holds project info to be saved.
type CreateProjectRequest struct {
	Title            string `json:"title"`
	OrganizationName string `json:"organizationName"`
	Description      string `json:"description"`
	PassType         string `json:"passType"`
}

// IsValid checks whether input is valid or not.
func (r *CreateProjectRequest) IsValid() error {
	if r.Title == "" {
		return errors.New("titlee is empty")
	}
	if r.OrganizationName == "" {
		return errors.New("organization name is empty")
	}
	if r.Description == "" {
		return errors.New("description is empty")
	}

	switch api.PassType(r.PassType) {
	case api.BoardingPass,
		api.Coupon,
		api.EventTicket,
		api.Generic,
		api.StoreCard:
		break
	default:
		return errors.New("pass type is invalid")
	}

	return nil
}

// String returns string representation of struct.
func (r *CreateProjectRequest) String() string {
	return fmt.Sprintf(
		`{"title":"%s","organizationName":"%s","description":"%s","passType":"%s"}`,
		r.Title,
		r.OrganizationName,
		r.Description,
		r.PassType,
	)
}

func (s *Service) createProjectHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req CreateProjectRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	project := api.NewProject(
		req.Title,
		req.OrganizationName,
		req.Description,
		api.PassType(req.PassType),
	)

	exists, err := s.env.Logic.IsProjectExists(
		ctx,
		project.Title,
		project.OrganizationName,
		project.Description,
		project.PassType,
	)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsProjectExists", err)
	}
	if exists {
		return sendJSON(w, http.StatusNotAcceptable, M{
			"title":            req.Title,
			"organizationName": req.OrganizationName,
			"description":      req.Description,
			"passType":         req.PassType,
		})
	}

	err = s.env.Logic.SaveNewProject(ctx, user, project)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SaveNewProject", err)
	}

	return sendJSON(w, http.StatusCreated, M{"id": project.ID})
}
