package redistore

import (
	"context"
	"errors"
	"time"
)

// InsertPass ...
func (p *Pool) InsertPass(ctx context.Context, serialNumber, authToken, passTypeID string) error {
	return errors.New("not implemented")
}

// UpdatePass ...
func (p *Pool) UpdatePass(ctx context.Context, serialNumber string) error {
	return errors.New("not implemented")
}

// FindPass ...
func (p *Pool) FindPass(ctx context.Context, serialNumber, authToken, passTypeID string) (bool, error) {
	return false, errors.New("not implemented")
}

// FindRegistration ...
func (p *Pool) FindRegistration(ctx context.Context, deviceID, serialNumber string) (bool, error) {
	return false, errors.New("not implemented")
}

// FindSerialNumbers ...
func (p *Pool) FindSerialNumbers(ctx context.Context, deviceID, passTypeID, tag string) ([]string, error) {
	return nil, errors.New("not implemented")
}

// LatestPass ...
func (p *Pool) LatestPass(ctx context.Context, serialNumber, authToken, passTypeID string) (time.Time, error) {
	return time.Time{}, errors.New("not implemented")
}

// InsertRegistration ...
func (p *Pool) InsertRegistration(ctx context.Context, deviceID, pushToken, serialNumber, passTypeID string) error {
	return errors.New("not implemented")
}

// DeleteRegistration ...
func (p *Pool) DeleteRegistration(ctx context.Context, deviceID, serialNumber, passTypeID string) (bool, error) {
	return false, errors.New("not implemented")
}

// InsertLog ...
func (p *Pool) InsertLog(ctx context.Context, remoteAddr, requestID, message string) error {
	return errors.New("not implemented")
}
