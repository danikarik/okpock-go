package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

const (
	backgroundImage = "background"
	footerImage     = "footer"
	iconImage       = "icon"
	logoImage       = "logo"
	stripImage      = "strip"
)

// UploadImageRequest holds image type and uuid from uploads.
type UploadImageRequest struct {
	UUID string        `json:"uuid"`
	Type string        `json:"type"`
	Size api.ImageSize `json:"size"`
}

// IsValid checks whether input is valid or not.
func (r *UploadImageRequest) IsValid() error {
	if r.UUID == "" {
		return errors.New("uuid is empty")
	}

	switch r.Type {
	case backgroundImage, footerImage, iconImage, logoImage, stripImage:
		break
	default:
		return errors.New("image type is invalid")
	}

	switch r.Size {
	case api.ImageSize1x, api.ImageSize2x, api.ImageSize3x:
		break
	default:
		return errors.New("image size is invalid")
	}

	return nil
}

// String returns string representation of struct.
func (r *UploadImageRequest) String() string {
	return fmt.Sprintf(
		`{"uuid":"%s","type":"%s","size":"%s"}`,
		r.UUID,
		r.Type,
		r.Size,
	)
}

func (s *Service) uploadProjectImage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req UploadImageRequest
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

	_, err = s.env.Logic.LoadUploadByUUID(ctx, user, req.UUID)
	if err != nil {
		return sendJSON(w, http.StatusNotAcceptable, M{"uuid": req.UUID})
	}

	switch req.Type {
	case backgroundImage:
		err = s.env.Logic.SetBackgroundImage(ctx, req.Size, req.UUID, project)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SetBackgroundImage", err)
		}
		break
	case footerImage:
		err = s.env.Logic.SetFooterImage(ctx, req.Size, req.UUID, project)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SetFooterImage", err)
		}
		break
	case iconImage:
		err = s.env.Logic.SetIconImage(ctx, req.Size, req.UUID, project)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SetIconImage", err)
		}
		break
	case logoImage:
		err = s.env.Logic.SetLogoImage(ctx, req.Size, req.UUID, project)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SetLogoImage", err)
		}
		break
	case stripImage:
		err = s.env.Logic.SetStripImage(ctx, req.Size, req.UUID, project)
		if err != nil {
			return s.httpError(w, r, http.StatusInternalServerError, "SetStripImage", err)
		}
		break
	}

	return sendJSON(w, http.StatusOK, project)
}
