package service

import (
	"net/http"
	"time"

	"github.com/danikarik/okpock/pkg/store/sequel"
	"github.com/danikarik/mux"
)

// Passes represents list of serial numbers
// associated with device id.
type Passes struct {
	LastUpdated   string   `json:"lastUpdated"`
	SerialNumbers []string `json:"serialNumbers"`
}

// SerialNumbers is used for
// "Getting the Serial Numbers for Passes Associated with a Device".
func (s *Service) serialNumbers(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx                     = r.Context()
		vars                    = mux.Vars(r)
		deviceLibraryIdentifier = vars["deviceLibraryIdentifier"]
		passTypeIdentifier      = vars["passTypeIdentifier"]
		passesUpdatedSince      = r.FormValue("passesUpdatedSince")
	)

	serials, err := s.env.PassKit.FindSerialNumbers(ctx, deviceLibraryIdentifier, passTypeIdentifier, passesUpdatedSince)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "FindSerialNumbers", err)
	}
	if len(serials) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	passes := &Passes{
		LastUpdated:   time.Now().Format(sequel.TimeFormat),
		SerialNumbers: serials,
	}

	err = sendJSON(w, http.StatusOK, passes)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Send", err)
	}

	return nil
}
