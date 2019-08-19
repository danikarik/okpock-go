package sequel

import (
	"context"
	"errors"

	"github.com/danikarik/okpock/pkg/api"
)

// IsOrganizationExists ...
func (m *MySQL) IsOrganizationExists(ctx context.Context, title string, userID int64) (bool, error) {
	return false, errors.New("not implemented")
}

// SaveNewOrganization ...
func (m *MySQL) SaveNewOrganization(ctx context.Context, org *api.Organization) error {
	return errors.New("not implemented")
}

// LoadOrganization ...
func (m *MySQL) LoadOrganization(ctx context.Context, id int64) (*api.Organization, error) {
	return nil, errors.New("not implemented")
}

// LoadOrganizations ...
func (m *MySQL) LoadOrganizations(ctx context.Context, userID int64) ([]*api.Organization, error) {
	return nil, errors.New("not implemented")
}

// UpdateOrganizationDescription ...
func (m *MySQL) UpdateOrganizationDescription(ctx context.Context, desc string, org *api.Organization) error {
	return errors.New("not implemented")
}

// UpdateOrganizationMetaData ...
func (m *MySQL) UpdateOrganizationMetaData(ctx context.Context, data map[string]interface{}, org *api.Organization) error {
	return errors.New("not implemented")
}
