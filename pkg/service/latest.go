package service

import (
	"net/http"
	"time"

	"github.com/danikarik/mux"
)

// LatestPass is used for
// "Getting the Latest Version of a Pass".
func (s *Service) latestPass(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx          = r.Context()
		vars         = mux.Vars(r)
		serialNumber = vars["serialNumber"]
		passTypeID   = vars["passTypeID"]
		authToken    = applePassFromContext(ctx)
	)

	lastUpdate, err := s.env.PassKit.LatestPass(ctx, serialNumber, authToken, passTypeID)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LatestPass", err)
	}

	if notModified(r, lastUpdate) {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	obj, err := s.env.Storage.GetFile(ctx, s.env.Config.PassesBucket, serialNumber)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "File", err)
	}

	w.Header().Set("Last-Modified", lastUpdate.Format(http.TimeFormat))
	err = obj.Serve(w)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Serve", err)
	}

	return nil
}

func notModified(r *http.Request, lastUpdate time.Time) bool {
	header := r.Header.Get("If-Modified-Since")
	if header == "" {
		return false
	}

	t, err := time.Parse(http.TimeFormat, header)
	if err != nil {
		return false
	}

	hasDiff := (lastUpdate.Year() > t.Year()) ||
		(lastUpdate.Month() > t.Month()) ||
		(lastUpdate.Day() > t.Day()) ||
		(lastUpdate.Hour() > t.Hour()) ||
		(lastUpdate.Minute() > t.Minute()) ||
		(lastUpdate.Second() > t.Second())

	return !hasDiff
}
