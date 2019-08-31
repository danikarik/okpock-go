package service

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/filestore"
	uuid "github.com/satori/go.uuid"
)

// MB is a megabyte.
const MB = 1 << 20

func (s *Service) readImageUpload(r *http.Request) (*api.Upload, error) {
	err := r.ParseMultipartForm(10 * MB)
	if err != nil {
		return nil, err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	hash, err := api.Hash(data)
	if err != nil {
		return nil, err
	}

	return &api.Upload{
		Filename:    header.Filename,
		Hash:        hash,
		Body:        data,
		ContentType: http.DetectContentType(data),
		CreatedAt:   time.Now(),
	}, nil
}

func (s *Service) uploadHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	upload, err := s.readImageUpload(r)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadImageUpload", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	exists, err := s.env.Logic.IsUploadExists(ctx, user, upload.Filename, upload.Hash)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "IsUploadExists", err)
	}

	if exists {
		return sendJSON(w, http.StatusNotAcceptable, upload)
	}

	bucket := s.env.Config.UploadBucket
	object := &filestore.Object{
		Prefix:      strconv.FormatInt(user.ID, 10),
		Key:         uuid.NewV4().String(),
		Body:        upload.Body,
		ContentType: upload.ContentType,
	}

	err = s.env.Storage.UploadFile(ctx, bucket, object)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UploadFile", err)
	}

	upload.UUID = object.Path()
	err = s.env.Logic.SaveNewUpload(ctx, user, upload)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "SaveNewUpload", err)
	}

	return sendJSON(w, http.StatusCreated, upload)
}
