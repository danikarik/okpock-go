package service

import (
	"net/http"

	"github.com/danikarik/mux"
)

// UnregisterDevice is used for
// "Unregistering a Device".
func (s *Service) unregisterDevice(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx          = r.Context()
		vars         = mux.Vars(r)
		deviceID     = vars["deviceID"]
		passTypeID   = vars["passTypeID"]
		serialNumber = vars["serialNumber"]
	)

	ok, err := s.env.PassKit.DeleteRegistration(ctx, deviceID, serialNumber, passTypeID)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "DeleteRegistration", err)
	}

	if !ok {
		return s.httpError(w, r, http.StatusNotFound, "", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
