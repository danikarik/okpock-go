package service

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/danikarik/mux"
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/secure"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	applePassContextKey contextKey = "apple_pass"
	requestIDKey        contextKey = "request_id"
	userContextKey      contextKey = "user"
)

var (
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")

	applePassRegexp = regexp.MustCompile(`^(?:A|a)pplePass (\S+$)`)
)

const csrfHeader string = "X-XSRF-TOKEN"

const (
	metaReferer               string = "referer"
	metaSuggestChangeUsername string = "suggestChangeUsername"
)

var safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}

// ErrMissingContext returned when context value is missing.
var ErrMissingContext = errors.New("context: missing value")

func withApplePass(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, applePassContextKey, token)
}

func applePassFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if token, ok := ctx.Value(applePassContextKey).(string); ok {
		return token
	}
	return ""
}

func parseAuthHeader(header string, re *regexp.Regexp) string {
	if header == "" {
		return ""
	}
	matches := re.FindStringSubmatch(header)
	if len(matches) != 2 {
		return ""
	}
	return matches[1]
}

func applePassMiddleware(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var (
		ctx  = r.Context()
		code = http.StatusUnauthorized
		err  = mux.NewHTTPError(code, http.StatusText(code))
	)
	token := parseAuthHeader(r.Header.Get("Authorization"), applePassRegexp)
	if token == "" {
		return ctx, err
	}
	return withApplePass(ctx, token), nil
}

func loggerMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		ctx := r.Context()
		start := time.Now()
		defer func() {
			logger.Info("served",
				zap.String("method", r.Method),
				zap.String("proto", r.Proto),
				zap.String("path", r.URL.Path),
				zap.Duration("lat", time.Since(start)),
				zap.String("reqId", reqIDFromContext(ctx)))
		}()
		return ctx, nil
	}
}

func withReqID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func reqIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func requestIDMiddleware(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = uuid.NewV4().String()
		r.Header.Set("X-Request-ID", requestID)
	}
	return withReqID(r.Context(), requestID), nil
}

func recovererMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logger.Error("recover", zap.ByteString("stack", debug.Stack()))
				code := http.StatusInternalServerError
				http.Error(w, http.StatusText(code), code)
			}
		}()
		return r.Context(), nil
	}
}

func realIPMiddleware(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	if rip := realIP(r); rip != "" {
		r.RemoteAddr = rip
	}
	return r.Context(), nil
}

func realIP(r *http.Request) string {
	var ip string

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	}

	return ip
}

func withUser(ctx context.Context, u *api.User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func userFromContext(ctx context.Context) (*api.User, error) {
	if u, ok := ctx.Value(userContextKey).(*api.User); ok {
		return u, nil
	}
	return nil, ErrMissingContext
}

func (s *Service) authMiddleware(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var (
		ctx  = r.Context()
		ucl  = NewClaims()
		code = http.StatusUnauthorized
	)

	err := s.getClaims(r, ucl)
	if err != nil {
		return nil, s.httpError(w, r, code, "GetClaims", err)
	}

	id, err := strconv.ParseInt(ucl.Subject, 10, 64)
	if err != nil {
		return nil, s.httpError(w, r, code, "ParseInt", err)
	}

	user, err := s.env.Auth.LoadUser(ctx, id)
	if err != nil {
		return nil, s.httpError(w, r, code, "LoadUser", err)
	}

	return withUser(ctx, user), nil
}

func newCSRFToken() string { return secure.Token() }

func skipCSRFCheck(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if strings.Contains(origin, "localhost") {
		return true
	}
	for _, method := range safeMethods {
		if method == r.Method {
			return true
		}
	}
	return false
}

func (s *Service) csrfMiddleware(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	var (
		ctx  = r.Context()
		ucl  = NewClaims()
		code = http.StatusForbidden
	)

	err := s.getClaims(r, ucl)
	if err != nil {
		return nil, s.httpError(w, r, code, "GetClaims", err)
	}

	if !skipCSRFCheck(r) {
		headerToken := r.Header.Get(csrfHeader)
		if headerToken == "" || headerToken != ucl.CSRFToken {
			return nil, s.httpError(w, r, code, "CSRFHeader", nil)
		}
	}

	return ctx, nil
}

func (s *Service) allowedOrigins() []string {
	origins := []string{
		"https://console.okpock.com",
	}
	if s.env.Config.IsDevelopment() {
		origins = []string{
			"https://console-dev.okpock.com",
		}
	}
	return origins
}

func (s *Service) corsMiddleware(next http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-Requested-With", csrfHeader},
		ExposedHeaders: []string{csrfHeader},
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			if s.env.Config.IsDevelopment() && strings.Contains(origin, "localhost") {
				return true
			}
			for _, allowedOrigin := range s.allowedOrigins() {
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		},
		AllowCredentials:   true,
		MaxAge:             300,
		OptionsPassthrough: false,
	})
	return cors.Handler(next)
}
