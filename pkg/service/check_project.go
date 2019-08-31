package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
)

// CheckProjectRequest holds project info to be checked.
type CheckProjectRequest struct {
	Title            string `json:"title"`
	OrganizationName string `json:"organizationName"`
	Description      string `json:"description"`
	PassType         string `json:"passType"`
}

// IsValid checks whether input is valid or not.
func (r *CheckProjectRequest) IsValid() error {
	if r.Title == "" {
		return errors.New("title is empty")
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
func (r *CheckProjectRequest) String() string {
	return fmt.Sprintf(
		`{"title":"%s","organizationName":"%s","description":"%s","passType":"%s"}`,
		r.Title,
		r.OrganizationName,
		r.Description,
		r.PassType,
	)
}

func (s *Service) checkProjectHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req CheckProjectRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	exists, err := s.env.Logic.IsProjectExists(
		ctx,
		req.Title,
		req.OrganizationName,
		req.Description,
		api.PassType(req.PassType),
	)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsProjectExists", err)
	}

	code := http.StatusOK
	if exists {
		code = http.StatusNotAcceptable
	}

	return sendJSON(w, code, M{
		"title":            req.Title,
		"organizationName": req.OrganizationName,
		"description":      req.Description,
		"passType":         req.PassType,
	})
}
