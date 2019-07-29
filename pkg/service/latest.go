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
		ctx                = r.Context()
		vars               = mux.Vars(r)
		serialNumber       = vars["serialNumber"]
		passTypeIdentifier = vars["passTypeIdentifier"]
		authToken          = applePassFromContext(ctx)
	)

	lastUpdate, err := s.env.PassKit.LatestPass(ctx, serialNumber, authToken, passTypeIdentifier)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "LatestPass", err)
	}

	if isModifiedSince(r, lastUpdate) {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	obj, err := s.env.Storage.GetFile(ctx, s.env.Config.PassesBucket, serialNumber+".pkpass")
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "File", err)
	}

	err = obj.Serve(w)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "Serve", err)
	}

	return nil
}

func isModifiedSince(r *http.Request, lastUpdate time.Time) bool {
	header := r.Header.Get("If-Modified-Since")
	if header == "" {
		return false
	}

	t, err := time.Parse(http.TimeFormat, header)
	if err != nil {
		return false
	}

	return lastUpdate.Before(t.Add(1 * time.Second))
}
