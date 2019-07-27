package service

import (
	"net/http"

	"go.uber.org/zap"
)

// ErrorLog represents log message from devices.
type ErrorLog struct {
	Logs []string `json:"logs"`
}

// IsValid checks whether input is valid or not.
func (e *ErrorLog) IsValid() error { return nil }

// String returns string representation of struct.
func (e *ErrorLog) String() string { return "" }

// ErrorLogs is used for
// "Logging Errors".
func (s *Service) errorLogs(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx   = r.Context()
		reqID = reqIDFromContext(ctx)
	)

	var errLog ErrorLog
	err := readJSON(r, &errLog)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "Read", err)
	}

	for _, msg := range errLog.Logs {
		err = s.env.PassKit.InsertLog(ctx, r.RemoteAddr, reqID, msg)
		if err != nil {
			s.logger.Error("InsertLog", zap.Error(err))
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
