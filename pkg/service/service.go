package service

import (
	"net/http"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/mux"
	"go.uber.org/zap"
)

// Service holds env and routes.
type Service struct {
	env     *env.Env
	logger  *zap.Logger
	handler http.Handler
}

// New returns a new instance of `Service`.
func New(env *env.Env, logger *zap.Logger) *Service {
	serverSigningSecret = []byte(env.Config.ServerSecret)

	srv := &Service{
		env:    env,
		logger: logger,
	}

	return srv.withRouter()
}

// ServeHTTP implemenents `http.Handler` interface.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Service) requestLogger(r *http.Request) *zap.Logger {
	return s.logger.
		With(zap.String("url", r.URL.String())).
		With(zap.String("method", r.Method))
}

func (s *Service) httpError(w http.ResponseWriter, r *http.Request, code int, msg string, err error) error {
	reqID := reqIDFromContext(r.Context())
	logger := s.requestLogger(r)
	if err != nil {
		logger.Error(
			"http_error",
			zap.Error(err),
			zap.Int("code", code),
			zap.String("message", msg),
			zap.String("request_id", reqID),
		)
	}
	return mux.NewHTTPError(code, http.StatusText(code)).
		WithErrorID(reqID).
		WithShowError(!s.env.Config.IsProduction()).
		WithInternalMessage(msg).
		WithInternalError(err)
}
