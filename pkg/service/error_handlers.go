package service

import (
	"net/http"

	"github.com/danikarik/mux"
	"go.uber.org/zap"
)

func (s *Service) errorHandler(err error, w http.ResponseWriter, r *http.Request) {
	logger := s.requestLogger(r)
	switch e := err.(type) {
	case *mux.HTTPError:
		if jsonErr := sendJSON(w, e.Code, e); jsonErr != nil {
			s.errorHandler(jsonErr, w, r)
		}
		break
	default:
		httpErr := mux.NewHTTPError(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		).
			WithErrorID(reqIDFromContext(r.Context())).
			WithInternalError(err)
		if jsonErr := sendJSON(w, httpErr.Code, httpErr); jsonErr != nil {
			logger.Error("error_handler", zap.Error(jsonErr))
		}
		break
	}
}

func (s *Service) notFoundHandler(w http.ResponseWriter, r *http.Request) error {
	return s.httpError(w, r, http.StatusNotFound, "404: url="+r.URL.String(), nil)
}

func (s *Service) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) error {
	return s.httpError(w, r, http.StatusMethodNotAllowed, "405: url="+r.URL.String()+", method="+r.Method, nil)
}
