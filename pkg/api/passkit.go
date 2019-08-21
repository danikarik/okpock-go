package api

import (
	"context"
	"time"
)

// PassKit implements `Apple PassKit` spec.
type PassKit interface {
	// InsertPass ...
	// TODO: description
	InsertPass(ctx context.Context, serialNumber, authToken, passTypeID string) error

	// UpdatePass ...
	// TODO: description
	UpdatePass(ctx context.Context, serialNumber string) error

	// FindPass ...
	// TODO: description
	FindPass(ctx context.Context, serialNumber, authToken, passTypeID string) (bool, error)

	// FindRegistration ...
	// TODO: description
	FindRegistration(ctx context.Context, deviceID, serialNumber string) (bool, error)

	// FindSerialNumbers ...
	// TODO: description
	FindSerialNumbers(ctx context.Context, deviceID, passTypeID, tag string) ([]string, error)

	// LatestPass ...
	// TODO: description
	LatestPass(ctx context.Context, serialNumber, authToken, passTypeID string) (time.Time, error)

	// InsertRegistration ...
	// TODO: description
	InsertRegistration(ctx context.Context, deviceID, pushToken, serialNumber, passTypeID string) error

	// DeleteRegistration ...
	// TODO: description
	DeleteRegistration(ctx context.Context, deviceID, serialNumber, passTypeID string) (bool, error)

	// InsertLog ...
	// TODO: description
	InsertLog(ctx context.Context, remoteAddr, requestID, message string) error
}
