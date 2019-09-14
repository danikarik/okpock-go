package service

import (
	"net/http"

	"github.com/danikarik/mux"
)

func (s *Service) downloadPkpass(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx          = r.Context()
		vars         = mux.Vars(r)
		serialNumber = vars["serialNumber"]
	)

	found, err := s.env.PassKit.FindRegistrationBySerialNumber(ctx, serialNumber)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "FindRegistrationBySerialNumber", err)
	}
	if found {
		return s.redirect(w, r, s.appURL(""))
	}

	obj, err := s.env.Storage.GetFile(ctx, s.env.Config.PassesBucket, serialNumber)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "File", err)
	}

	err = obj.Serve(w)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Serve", err)
	}

	return nil
}
