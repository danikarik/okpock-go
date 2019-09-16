package api

import (
	"context"
	"time"
)

// PassKit implements `Apple PassKit` spec.
type PassKit interface {
	// InsertPass ...
	InsertPass(ctx context.Context, serialNumber, authToken, passTypeID string) error

	// UpdatePass ...
	UpdatePass(ctx context.Context, serialNumber string) error

	// FindPass ...
	FindPass(ctx context.Context, serialNumber, authToken, passTypeID string) (bool, error)

	// FindRegistration ...
	FindRegistration(ctx context.Context, deviceID, serialNumber string) (bool, error)

	// FindRegistrationBySerialNumber ...
	FindRegistrationBySerialNumber(ctx context.Context, serialNumber string) (bool, error)

	// FindPushToken
	FindPushToken(ctx context.Context, serialNumber string) (string, error)

	// FindSerialNumbers ...
	FindSerialNumbers(ctx context.Context, deviceID, passTypeID, tag string) ([]string, error)

	// LatestPass ...
	LatestPass(ctx context.Context, serialNumber, authToken, passTypeID string) (time.Time, error)

	// InsertRegistration ...
	InsertRegistration(ctx context.Context, deviceID, pushToken, serialNumber, passTypeID string) error

	// DeleteRegistration ...
	DeleteRegistration(ctx context.Context, deviceID, serialNumber, passTypeID string) (bool, error)

	// InsertLog ...
	InsertLog(ctx context.Context, remoteAddr, requestID, message string) error
}
