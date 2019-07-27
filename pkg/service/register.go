package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danikarik/mux"
)

// Register represents register endpoint's payload.
// It holds device's push token.
type Register struct {
	PushToken string `json:"pushToken"`
}

// IsValid checks whether input is valid or not.
func (r *Register) IsValid() error {
	if r.PushToken == "" {
		return errors.New("push token is empty")
	}
	return nil
}

// String returns string representation of struct.
func (r *Register) String() string { return fmt.Sprintf(`{"pushToken":"%s"}`, r.PushToken) }

// RegisterDevice is used for
// "Registering a Device to Receive Push Notifications for a Pass".
func (s *Service) registerDevice(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx                     = r.Context()
		vars                    = mux.Vars(r)
		deviceLibraryIdentifier = vars["deviceLibraryIdentifier"]
		passTypeIdentifier      = vars["passTypeIdentifier"]
		serialNumber            = vars["serialNumber"]
	)

	var register Register
	err := readJSON(r, &register)
	if err != nil {
		return s.httpError(w, r, http.StatusBadRequest, "Read", err)
	}

	exists, err := s.env.PassKit.FindRegistration(ctx, deviceLibraryIdentifier, serialNumber)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "FindRegistration", err)
	}
	if exists {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	err = s.env.PassKit.InsertRegistration(ctx, deviceLibraryIdentifier, register.PushToken, serialNumber, passTypeIdentifier)
	if err != nil {
		return s.httpError(w, r, http.StatusInternalServerError, "InsertRegistration", err)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}
