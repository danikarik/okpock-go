package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/danikarik/okpock/pkg/api"
)

// MetaDataChangeRequest holds new user data.
type MetaDataChangeRequest struct {
	Data api.JSONMap `json:"data"`
}

// IsValid checks whether input is valid or not.
func (r *MetaDataChangeRequest) IsValid() error {
	if r.Data == nil || len(r.Data) == 0 {
		return errors.New("data is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *MetaDataChangeRequest) String() string {
	json, err := json.Marshal(r.Data)
	if err != nil {
		return "{}"
	}
	return string(json)
}

func (s *Service) metaDataChangeHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req MetaDataChangeRequest
	err := readJSON(r, &req)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "ReadJSON", err)
	}

	user, err := userFromContext(ctx)
	if err != nil {
		return s.httpError(w, r, http.StatusUnauthorized, "UserFromContext", err)
	}

	err = s.env.Auth.UpdateUserMetaData(ctx, req.Data, user)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "UpdateUserMetaData", err)
	}

	return sendJSON(w, http.StatusOK, M{
		"id":           user.ID,
		"userMetadata": user.UserMetaData,
		"updatedAt":    user.UpdatedAt,
	})
}
