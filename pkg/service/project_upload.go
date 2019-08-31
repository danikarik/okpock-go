package service

import (
	"io/ioutil"
	"net/http"
)

// MB is a megabyte.
const MB = 1 << 20

const (
	backgroundImage = "background"
	footerImage     = "footer"
	iconImage       = "icon"
	stripImage      = "strip"
)

func (s *Service) readImageUpload(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(10 * MB)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func (s *Service) uploadProjectImage(w http.ResponseWriter, r *http.Request) error {
	// ctx := r.Context()

	// data, err := s.readImageUpload(r)
	// if err != nil {
	// 	return s.httpError(w, r, http.StatusBadRequest, "ReadImageUpload", err)
	// }

	// user, err := userFromContext(ctx)
	// if err != nil {
	// 	return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	// }

	// id, err := s.idFromRequest(r, "id")
	// if err != nil {
	// 	return s.httpError(w, r, http.StatusBadRequest, "IDFromRequest", err)
	// }

	// project, err := s.env.Logic.LoadProject(ctx, user, id)
	// if err == store.ErrNotFound {
	// 	return s.httpError(w, r, http.StatusNotFound, "LoadProject", err)
	// }
	// if err != nil {
	// 	return s.httpError(w, r, http.StatusInternalServerError, "LoadProject", err)
	// }

	// s.env.Storage.UploadFile(ctx, s.env.Config.bu)

	// TODO: SetBackgroundImage, SetFooterImage, SetIconImage, SetStripImage
	return sendJSON(w, 200, M{})
}
