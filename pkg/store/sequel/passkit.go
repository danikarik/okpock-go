package sequel

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/store"
)

// InsertPass ...
func (m *MySQL) InsertPass(ctx context.Context, serialNumber, authToken, passTypeID string) error {
	query := m.builder.Insert("passes").
		Columns("serial_number", "authentication_token", "pass_type_id", "updated_at").
		Values(serialNumber, authToken, passTypeID, time.Now())

	_, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePass ...
func (m *MySQL) UpdatePass(ctx context.Context, serialNumber string) error {
	query := m.builder.Update("passes").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"serial_number": serialNumber})

	_, err := m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// FindPass ...
func (m *MySQL) FindPass(ctx context.Context, serialNumber, authToken, passTypeID string) (bool, error) {
	query := m.builder.Select("count(1)").From("passes").
		Where(sq.Eq{
			"serial_number":        serialNumber,
			"authentication_token": authToken,
			"pass_type_id":         passTypeID,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// FindRegistration ...
func (m *MySQL) FindRegistration(ctx context.Context, deviceID, serialNumber string) (bool, error) {
	query := m.builder.Select("count(1)").From("registrations").
		Where(sq.Eq{
			"device_id":     deviceID,
			"serial_number": serialNumber,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// FindRegistrationBySerialNumber ...
func (m *MySQL) FindRegistrationBySerialNumber(ctx context.Context, serialNumber string) (bool, error) {
	query := m.builder.Select("count(1)").From("registrations").
		Where(sq.Eq{
			"serial_number": serialNumber,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// FindPushToken ...
func (m *MySQL) FindPushToken(ctx context.Context, serialNumber string) (string, error) {
	query := m.builder.Select("push_token").From("registrations").
		Where(sq.Eq{
			"serial_number": serialNumber,
		})

	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return "", err
	}

	var pushToken string
	err = row.Scan(&pushToken)
	if err == sql.ErrNoRows {
		return "", store.ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return pushToken, nil
}

// FindSerialNumbers ...
func (m *MySQL) FindSerialNumbers(ctx context.Context, deviceID, passTypeID, tag string) ([]string, error) {
	var sns []string

	query := m.builder.Select("p.serial_number").From("passes p").
		LeftJoin("registrations r on r.serial_number = p.serial_number").
		Where(sq.Eq{
			"r.device_id":    deviceID,
			"r.pass_type_id": passTypeID,
		})

	if tag != "" {
		query = query.Where(sq.Gt{"p.updated_at": tag})
	}

	rows, err := m.selectQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var sn string
		if err := rows.Scan(&sn); err != nil {
			return nil, err
		}
		sns = append(sns, sn)
	}

	return sns, nil
}

// LatestPass ...
func (m *MySQL) LatestPass(ctx context.Context, serialNumber, authToken, passTypeID string) (time.Time, error) {
	var t time.Time

	query := m.builder.Select("updated_at").From("passes").
		Where(sq.Eq{
			"serial_number":        serialNumber,
			"authentication_token": authToken,
			"pass_type_id":         passTypeID,
		})

	err := m.scanQuery(ctx, query, &t)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// InsertRegistration ...
func (m *MySQL) InsertRegistration(ctx context.Context, deviceID, pushToken, serialNumber, passTypeID string) error {
	query := m.builder.Insert("registrations").
		Columns("device_id", "push_token", "serial_number", "pass_type_id").
		Values(deviceID, pushToken, serialNumber, passTypeID)

	_, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRegistration ...
func (m *MySQL) DeleteRegistration(ctx context.Context, deviceID, serialNumber, passTypeID string) (bool, error) {
	query := m.builder.Delete("registrations").
		Where(sq.Eq{
			"device_id":     deviceID,
			"serial_number": serialNumber,
			"pass_type_id":  passTypeID,
		})

	rows, err := m.deleteQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// InsertLog ...
func (m *MySQL) InsertLog(ctx context.Context, remoteAddr, requestID, message string) error {
	query := m.builder.Insert("logs").
		Columns("remote_address", "request_id", "message", "updated_at").
		Values(remoteAddr, requestID, message, time.Now())

	_, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
