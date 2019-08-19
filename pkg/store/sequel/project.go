package sequel

import (
	"context"
	"errors"

	"github.com/danikarik/okpock/pkg/api"
)

// IsProjectExists ...
func (m *MySQL) IsProjectExists(ctx context.Context, desc string, orgID int64, passType api.PassType) (bool, error) {
	return false, errors.New("not implemented")
}

// SaveNewProject ...
func (m *MySQL) SaveNewProject(ctx context.Context, proj *api.Project) error {
	return errors.New("not implemented")
}

// LoadProject ...
func (m *MySQL) LoadProject(ctx context.Context, id int64) (*api.Project, error) {
	return nil, errors.New("not implemented")
}

// LoadProjects ...
func (m *MySQL) LoadProjects(ctx context.Context, userID int64) ([]*api.Project, error) {
	return nil, errors.New("not implemented")
}

// UpdateProjectDescription ...
func (m *MySQL) UpdateProjectDescription(ctx context.Context, desc string, proj *api.Project) error {
	return errors.New("not implemented")
}

// SetBackgroundImage ...
func (m *MySQL) SetBackgroundImage(ctx context.Context, key string, proj *api.Project) error {
	return errors.New("not implemented")
}

// SetFooterImage ...
func (m *MySQL) SetFooterImage(ctx context.Context, key string, proj *api.Project) error {
	return errors.New("not implemented")
}

// SetIconImage ...
func (m *MySQL) SetIconImage(ctx context.Context, key string, proj *api.Project) error {
	return errors.New("not implemented")
}

// SetStripImage ...
func (m *MySQL) SetStripImage(ctx context.Context, key string, proj *api.Project) error {
	return errors.New("not implemented")
}
